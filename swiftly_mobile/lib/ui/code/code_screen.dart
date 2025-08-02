import 'package:flutter/material.dart';
import 'dart:ui';
import 'widgets/explorer_sidebar.dart';
import 'widgets/editor_tabbar.dart';
import 'widgets/main_editor_area.dart';

final ValueNotifier<String> activeTab = ValueNotifier("CameraViewModel.swift");
final ValueNotifier<List<String>> openTabs = ValueNotifier(["CameraViewModel.swift"]);

class CodeScreen extends StatelessWidget {
  const CodeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F172A),
      body: Row(
        children: [
          const ExplorerSidebar(),
          Expanded(
            child: Column(
              children: [
                const EditorTabBar(),
                Expanded(
                  child: Container(
                    color: const Color(0xFF1E293B),
                    child: ValueListenableBuilder<String>(
                      valueListenable: activeTab,
                      builder: (context, file, _) {
                        final content = '''
                          import UIKit

                          class CameraViewModel {
                              let userDefaults = UserDefaults.standard
                              let dbManager = DBService()

                              var isFlashEnabled: Bool {
                                  get {
                                      return userDefaults.bool(forKey: "IsFlashEnabled")
                                  }
                                  set {
                                      userDefaults.set(newValue, forKey: "IsFlashEnabled")
                                  }
                              }
                          }
                          ''';
                        final contentLines = content.split('\n');
                        final lines = List.generate(
                          contentLines.length + 1,
                          (index) => '${index + 1}'.padLeft(3),
                        );
                        return Row(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Container(
                              width: 40,
                              color: const Color(0xFF1E293B),
                              padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 4),
                              child: ListView.builder(
                                itemCount: lines.length,
                                itemBuilder: (context, index) {
                                  return Text(
                                    lines[index],
                                    textAlign: TextAlign.right,
                                    style: const TextStyle(
                                      color: Colors.white38,
                                      fontSize: 13,
                                      height: 1.5,
                                      fontFamily: 'SourceCodePro',
                                    ),
                                  );
                                },
                              ),
                            ),
                            const VerticalDivider(width: 1, color: Color(0xFF334155)),
                            Expanded(
                              child: SingleChildScrollView(
                                padding: const EdgeInsets.all(16),
                                child: SelectableText(
                                  content,
                                  style: const TextStyle(
                                    fontFamily: 'SourceCodePro',
                                    fontSize: 14,
                                    color: Color(0xFFF1F5F9),
                                    height: 1.5,
                                  ),
                                ),
                              ),
                            ),
                          ],
                        );
                      },
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
