import 'package:flutter/material.dart';
import 'explorer_sidebar.dart';

class EditorTabBarTop extends StatelessWidget {
  const EditorTabBarTop({super.key});

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<List<String>>(
      valueListenable: openTabs,
      builder: (context, tabs, _) {
        return ValueListenableBuilder<String>(
          valueListenable: activeTab,
          builder: (context, current, __) {
            return Container(
              height: 58,
              decoration: BoxDecoration(
                color: const Color.fromRGBO(2, 10, 23, 1),
                border: Border(bottom: BorderSide(color: Colors.grey.shade300)),
              ),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  ...tabs.map((tab) => Container(
                        child: _EditorTab(
                          title: tab,
                          isActive: tab == current,
                          onClose: () {
                            final updated = List<String>.from(tabs)..remove(tab);
                            openTabs.value = updated;
                            if (current == tab && updated.isNotEmpty) {
                              activeTab.value = updated.last;
                            } else if (updated.isEmpty) {
                              activeTab.value = "";
                            }
                          },
                          onSelect: () => activeTab.value = tab,
                        ),
                      )),
                  const Spacer(),
                ],
              ),
            );
          },
        );
      },
    );
  }
}

class _EditorTab extends StatelessWidget {
  final String title;
  final bool isActive;
  final VoidCallback onClose;
  final VoidCallback onSelect;

  const _EditorTab({
    required this.title,
    required this.isActive,
    required this.onClose,
    required this.onSelect,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onSelect,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: isActive ? Color.fromRGBO(37, 58, 102, 1) : Color.fromRGBO(37, 58, 102, 0.3),
          borderRadius: BorderRadius.only(
            topLeft: Radius.circular(12),
            topRight: Radius.circular(12)
          ),
        ),
        child: Row(
          children: [
            Icon(Icons.description, size: 16, color: isActive ? Color.fromRGBO(76, 111, 185, 1) : Color.fromRGBO(179, 179, 179, 1)),
            const SizedBox(width: 6),
            Text(
              title,
              style: TextStyle(
                fontSize: 13,
                color: Color.fromRGBO(179, 179, 179, 1),
                fontWeight: FontWeight.w500,
              ),
            ),
            const SizedBox(width: 6),
            GestureDetector(
              onTap: () {
                onClose();
              },
              child: Icon(Icons.close, size: 16, color: Color.fromRGBO(179, 179, 179, 1)),
            ),
          ],
        ),
      ),
    );
  }
}

class BottomEditorToolbar extends StatelessWidget {
  const BottomEditorToolbar({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 34,
      decoration: BoxDecoration(
        color: const Color.fromRGBO(37, 58, 102, 0.25), // тёмно-синий фон
        border: Border(
          top: BorderSide(color: Color.fromRGBO(91, 92, 94, 1), width: 0.1),
        ),
      ),
      child: Row(
        children: [
          _ToolbarButton(icon: Icons.add, label: 'Code'),
          _Separator(),
          _ToolbarButton(icon: Icons.play_arrow, label: 'Run All'),
          _ToolbarButton(icon: Icons.restart_alt, label: 'Restart'),
          _ToolbarButton(icon: Icons.clear_all, label: 'Clear all outputs'),
          _Separator(),
          _ToolbarButton(icon: Icons.more_vert, label: ''),
          _ToolbarButton(icon: Icons.data_object, label: 'Variables'),
          _ToolbarButton(icon: Icons.segment, label: 'Outline'),
          _ToolbarButton(icon: Icons.more_horiz, label: ''),
        ],
      ),
    );
  }
}

class _ToolbarButton extends StatelessWidget {
  final IconData icon;
  final String label;

  const _ToolbarButton({required this.icon, required this.label});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 10),
      child: TextButton.icon(
        onPressed: () {},
        icon: Icon(icon, size: 18, color: Color.fromRGBO(179, 179, 179, 1)),
        label: Text(label, style: const TextStyle(color: Color.fromRGBO(179, 179, 179, 1))),
        style: TextButton.styleFrom(
          padding: const EdgeInsets.symmetric(horizontal: 8),
          foregroundColor: Colors.white,
        ),
      ),
    );
  }
}

class _Separator extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      height: 24,
      width: 1,
      color: Colors.white24,
      margin: const EdgeInsets.symmetric(horizontal: 6),
    );
  }
}