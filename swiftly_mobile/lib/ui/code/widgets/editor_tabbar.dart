import 'package:flutter/material.dart';
import 'explorer_sidebar.dart';

class EditorTabBar extends StatelessWidget {
  const EditorTabBar({super.key});

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<List<String>>(
      valueListenable: openTabs,
      builder: (context, tabs, _) {
        return ValueListenableBuilder<String>(
          valueListenable: activeTab,
          builder: (context, current, __) {
            return Container(
              height: 44,
              decoration: BoxDecoration(
                color: const Color(0xFF0F172A),
                border: Border(bottom: BorderSide(color: Colors.grey.shade300)),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.05),
                    blurRadius: 6,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              child: Row(
                children: [
                  const SizedBox(width: 16),
                  ...tabs.map((tab) => Padding(
                        padding: const EdgeInsets.only(right: 8),
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
                  IconButton(
                    icon: const Icon(Icons.play_arrow_rounded, size: 20),
                    tooltip: 'Run',
                    onPressed: () {},
                  ),
                  IconButton(
                    icon: const Icon(Icons.refresh_rounded, size: 20),
                    tooltip: 'Restart',
                    onPressed: () {},
                  ),
                  IconButton(
                    icon: const Icon(Icons.clear_all_rounded, size: 20),
                    tooltip: 'Clear',
                    onPressed: () {},
                  ),
                  const SizedBox(width: 8),
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
          color: isActive ? Colors.white : Colors.transparent,
          borderRadius: BorderRadius.circular(8),
          border: isActive
              ? Border.all(color: Colors.grey.shade300)
              : Border.all(color: Colors.transparent),
        ),
        child: Row(
          children: [
            Icon(Icons.description, size: 16, color: isActive ? Colors.black : Colors.black54),
            const SizedBox(width: 6),
            Text(
              title,
              style: TextStyle(
                fontSize: 13,
                color: isActive ? Colors.black : Colors.black54,
                fontWeight: FontWeight.w500,
              ),
            ),
            const SizedBox(width: 6),
            GestureDetector(
              onTap: () {
                onClose();
              },
              child: const Icon(Icons.close, size: 16, color: Colors.black38),
            ),
          ],
        ),
      ),
    );
  }
}
