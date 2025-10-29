#include "flutter_window.h"

#include <optional>
#include <dwmapi.h>
#pragma comment(lib, "dwmapi.lib")

#include "flutter/generated_plugin_registrant.h"

inline int GET_X_LPARAM(LPARAM lParam) {
  return (int)(short)LOWORD(lParam);
}

inline int GET_Y_LPARAM(LPARAM lParam) {
  return (int)(short)HIWORD(lParam);
}

LRESULT HitTestNCA(HWND hWnd, WPARAM wParam, LPARAM lParam) {
  POINT ptMouse = {GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam)};

  RECT rcWindow;
  GetWindowRect(hWnd, &rcWindow);

  const int BORDER_WIDTH = 8;
  const int TITLE_BAR_HEIGHT = 32;
  
  USHORT uRow = 1;
  USHORT uCol = 1;

  if (ptMouse.y >= rcWindow.top && ptMouse.y < rcWindow.top + TITLE_BAR_HEIGHT) {
    uRow = 0;
  }
  else if (ptMouse.y < rcWindow.bottom && ptMouse.y >= rcWindow.bottom - BORDER_WIDTH) {
    uRow = 2;
  }

  if (ptMouse.x >= rcWindow.left && ptMouse.x < rcWindow.left + BORDER_WIDTH) {
    uCol = 0;
  }
  else if (ptMouse.x < rcWindow.right && ptMouse.x >= rcWindow.right - BORDER_WIDTH) {
    uCol = 2;
  }

  LRESULT hitTests[3][3] = {
      {HTTOPLEFT, HTTOP, HTTOPRIGHT},
      {HTLEFT, HTNOWHERE, HTRIGHT},
      {HTBOTTOMLEFT, HTBOTTOM, HTBOTTOMRIGHT}
  };

  return hitTests[uRow][uCol];
}

FlutterWindow::FlutterWindow(const flutter::DartProject& project)
    : project_(project) {}

FlutterWindow::~FlutterWindow() {}

bool FlutterWindow::OnCreate() {
  if (!Win32Window::OnCreate()) {
    return false;
  }

  RECT frame = GetClientArea();

  flutter_controller_ = std::make_unique<flutter::FlutterViewController>(
      frame.right - frame.left, frame.bottom - frame.top, project_);
  if (!flutter_controller_->engine() || !flutter_controller_->view()) {
    return false;
  }
  RegisterPlugins(flutter_controller_->engine());
  SetChildContent(flutter_controller_->view()->GetNativeWindow());

  flutter_controller_->engine()->SetNextFrameCallback([&]() {
    this->Show();
  });

  flutter_controller_->ForceRedraw();

  return true;
}

void FlutterWindow::OnDestroy() {
  if (flutter_controller_) {
    flutter_controller_ = nullptr;
  }

  Win32Window::OnDestroy();
}

LRESULT
FlutterWindow::MessageHandler(HWND hwnd, UINT const message,
                              WPARAM const wparam,
                              LPARAM const lparam) noexcept {
  LRESULT lRet = 0;
  if (DwmDefWindowProc(hwnd, message, wparam, lparam, &lRet)) {
    return lRet;
  }

  if (flutter_controller_) {
    std::optional<LRESULT> result =
        flutter_controller_->HandleTopLevelWindowProc(hwnd, message, wparam,
                                                      lparam);
    if (result) {
      return *result;
    }
  }

  switch (message) {
    case WM_FONTCHANGE:
      flutter_controller_->engine()->ReloadSystemFonts();
      break;

    case WM_NCCALCSIZE: {
      if (wparam) {
        // Расширяем клиентскую область на весь размер окна
        // Это убирает видимую рамку DWM
        NCCALCSIZE_PARAMS* pncsp = reinterpret_cast<NCCALCSIZE_PARAMS*>(lparam);
        
        // Просто расширяем область на весь размер окна
        // Оставляем прямоугольник без изменений
        
        return WVR_REDRAW;
      }
      break;
    }

    case WM_NCHITTEST: {
      LRESULT hit = HitTestNCA(hwnd, wparam, lparam);
      if (hit != HTNOWHERE) {
        return hit;
      }
      break;
    }
  }

  return Win32Window::MessageHandler(hwnd, message, wparam, lparam);
}

