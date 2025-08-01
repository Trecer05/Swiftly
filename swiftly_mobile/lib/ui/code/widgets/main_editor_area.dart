import 'package:flutter/material.dart';
import 'explorer_sidebar.dart';

class MainEditorArea extends StatelessWidget {
  const MainEditorArea({super.key});

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<String>(
      valueListenable: activeTab,
      builder: (context, file, _) {
        return Container(
          color: const Color(0xFF0F172A),
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          child: Container(
            decoration: BoxDecoration(
              color: const Color(0xFF1E293B),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: const Color(0xFF334155)),
            ),
            padding: const EdgeInsets.all(16),
            child: file.isEmpty
                ? const Center(
                    child: Text(
                      "Нет открытых файлов",
                      style: TextStyle(color: Colors.white70),
                    ),
                  )
                : SingleChildScrollView(
                    child: SelectableText(
                      "Заглушка: содержимое файла \"$file\" будет отображаться здесь.",
                      style: const TextStyle(
                        fontFamily: 'SourceCodePro',
                        fontSize: 14,
                        color: Color(0xFFF1F5F9),
                        height: 1.5,
                      ),
                    ),
                  ),
          ),
        );
      },
    );
  }
}
