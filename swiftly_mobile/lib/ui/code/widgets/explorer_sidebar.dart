import 'package:flutter/material.dart';

final ValueNotifier<List<String>> openTabs = ValueNotifier(["CameraViewModel.swift"]);
final ValueNotifier<String> activeTab = ValueNotifier("CameraViewModel.swift");

class ExplorerSidebar extends StatefulWidget {
  const ExplorerSidebar({super.key});

  @override
  State<ExplorerSidebar> createState() => _ExplorerSidebarState();
}

class _ExplorerSidebarState extends State<ExplorerSidebar> {
  final Map<String, List<String>> directories = {
    'Html.Css_project': [
      'hub.html',
      'main.dart',
    ],
    'Docs': [
      'README.md',
    ],
  };

  final Set<String> expanded = {'Html.Css_project'};

  void toggleExpand(String dir) {
    setState(() {
      if (expanded.contains(dir)) {
        expanded.remove(dir);
      } else {
        expanded.add(dir);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 250,
      decoration: BoxDecoration(
        color: const Color(0xFF0F172A), // deep blue-gray from second screenshot
        border: const Border(
          right: BorderSide(color: Color(0xFF1E293B)),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            child: const Text(
              "Проводник",
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
                color: Color(0xFFF8FAFC),
              ),
            ),
          ),
          const Divider(height: 1, color: Color(0xFF1E293B)),
          Expanded(
            child: ListView(
              children: directories.keys.map((dir) {
                final isOpen = expanded.contains(dir);
                return Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    ListTile(
                      leading: Icon(
                        isOpen ? Icons.folder_open : Icons.folder,
                        color: const Color(0xFF38BDF8),
                        size: 20,
                      ),
                      title: Text(
                        dir,
                        style: const TextStyle(
                          color: Colors.white,
                          fontWeight: FontWeight.w500,
                          fontSize: 14,
                        ),
                      ),
                      onTap: () => toggleExpand(dir),
                    ),
                    if (isOpen)
                      Padding(
                        padding: const EdgeInsets.only(left: 16),
                        child: Column(
                          children: directories[dir]!
                              .map(
                                (file) => ListTile(
                                  dense: true,
                                  contentPadding: const EdgeInsets.only(left: 24, right: 8),
                                  leading: const Icon(Icons.insert_drive_file, size: 16, color: Colors.white70),
                                  title: Text(
                                    file,
                                    style: const TextStyle(
                                      fontSize: 13,
                                      color: Colors.white70,
                                    ),
                                  ),
                                  onTap: () {
                                    if (!openTabs.value.contains(file)) {
                                      openTabs.value = [...openTabs.value, file];
                                    }
                                    activeTab.value = file;
                                  },
                                ),
                              )
                              .toList(),
                        ),
                      ),
                  ],
                );
              }).toList(),
            ),
          ),
        ],
      ),
    );
  }
}
