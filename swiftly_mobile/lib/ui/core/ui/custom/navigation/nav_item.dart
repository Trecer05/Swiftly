import 'package:flutter/material.dart';

import '../../../../../routing/routers.dart';

enum NavItem {
  home(Routers.home, Icons.home, 'Главная'),
  chat(Routers.chat, Icons.chat, 'Чат'),
  code(Routers.code, Icons.code, 'Код'),
  cloud(Routers.cloud, Icons.cloud, 'Облако'),
  figma(Routers.figma, Icons.pan_tool, 'Фигма'),
  board(Routers.board, Icons.task, 'Задачи'),
  settings(Routers.settings, Icons.settings, 'Настройки');

  final String route;
  final IconData icon;
  final String label;

  const NavItem(this.route, this.icon, this.label);

  static final _routeMap = {for (var item in NavItem.values) item.route: item};

  static NavItem? fromRoute(String route) => _routeMap[route];
}