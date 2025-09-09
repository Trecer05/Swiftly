import 'dart:io' as io show File; // избегаем конфликта с Web

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:file_picker/file_picker.dart';

import 'package:swiftly_mobile/ui/cloud/widgets/left_panel.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/custom_app_bar_desktop.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/custom_button.dart';

class CloudScreen extends StatefulWidget {
  const CloudScreen({super.key});

  @override
  State<CloudScreen> createState() => _CloudScreenState();
}

class _CloudScreenState extends State<CloudScreen> {
  // Для десктопа/мобилок
  io.File? pickedFile;

  // Для Web
  Uint8List? pickedFileBytes;

  String? fileName;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color.fromARGB(255, 9, 30, 114),
      body: Row(
        children: [
          const SizedBox(width: 200, child: LeftPanel()),
          Expanded(
            child: Column(
              children: [
                CustomAppBarDesktop(title: 'Файлы', quantity: 16, buttons: [CustomButton(prefixIcon: Icons.search, gradient: false, onTap: (){}), CustomButton(prefixIcon: Icons.sort, text: 'По названию', suffixIcon: Icons.keyboard_arrow_down, gradient: false, onTap: (){}), CustomButton(prefixIcon: Icons.add, text: 'Добавить', gradient: true, onTap: (){})]),
                Expanded(
                  child: Container(
                    color: Colors.white,
                    child: Center(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          if (pickedFile != null && !kIsWeb)
                            SizedBox(
                              width: 250,
                              height: 250,
                              child: Image.file(pickedFile!),
                            ),
                          if (pickedFileBytes != null && kIsWeb)
                            SizedBox(
                              width: 250,
                              height: 250,
                              child: Image.memory(pickedFileBytes!),
                            ),
                          if (fileName != null) ...[
                            const SizedBox(height: 8),
                            Text(fileName!, style: const TextStyle(fontSize: 16)),
                          ],
                          if (pickedFile == null && pickedFileBytes == null)
                            ElevatedButton(
                              onPressed: _pickFile,
                              child: const Text('Добавить файл'),
                            ),
                        ],
                      ),
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

  Future<void> _pickFile() async {
    FilePickerResult? result = await FilePicker.platform.pickFiles(
      type: FileType.image,
    );

    if (result != null) {
      setState(() {
        fileName = result.files.single.name;

        if (kIsWeb) {
          pickedFileBytes = result.files.single.bytes;
        } else {
          pickedFile = io.File(result.files.single.path!);
        }
      });
    }
  }
}
