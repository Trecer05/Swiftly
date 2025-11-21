import 'package:flutter/material.dart';
import 'dart:io';
import 'widgets/file_model.dart';

class FileEditorScreen extends StatefulWidget {
  final FileInfo file;

  const FileEditorScreen({required this.file, super.key});

  @override
  State<FileEditorScreen> createState() => FileEditorScreenState();
}

class FileEditorScreenState extends State<FileEditorScreen> {
  late TextEditingController controller;
  bool isModified = false;

  @override
  void initState() {
    super.initState();
    controller = TextEditingController();
    loadFileContent();
  }

  Future<void> loadFileContent() async {
    if (widget.file.localPath == null) return;

    try {
      final file = File(widget.file.localPath!);
      if (widget.file.type == 'image') {
        setState(() {
          // Изображение будет показано через Image.file()
        });
      } else if (widget.file.type == 'file' && widget.file.name.endsWith('.txt')) {
        final content = await file.readAsString();
        controller.text = content;
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Ошибка: $e')),
      );
    }
  }

  Future<void> saveFile() async {
    if (widget.file.localPath == null) return;

    try {
      final file = File(widget.file.localPath!);
      await file.writeAsString(controller.text);
      setState(() {
        isModified = false;
      });
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Файл сохранён')),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Ошибка: $e')),
      );
    }
  }

  Widget _buildContent() {
    if (widget.file.type == 'image' && widget.file.localPath != null) {
      return Image.file(
        File(widget.file.localPath!),
        fit: BoxFit.contain,
      );
    } else if (widget.file.type == 'file' && widget.file.name.endsWith('.txt')) {
      return TextField(
        controller: controller,
        maxLines: null,
        expands: true,
        onChanged: (value) {
          setState(() {
            isModified = true;
          });
        },
        decoration: InputDecoration(
          border: InputBorder.none,
          filled: true,
          fillColor: Colors.transparent,
          contentPadding: EdgeInsets.all(12),
        ),
        style: TextStyle(color: Colors.white),
      );
    } else {
      return Center(
        child: Text(
          'Файл не поддерживается',
          style: TextStyle(color: Colors.white),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.transparent,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        leading: IconButton(
          icon: Icon(
            Icons.arrow_back,
            color: Color(0x80FFFFFF),
          ),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: Text(
          widget.file.name,
          style: const TextStyle(color: Color(0x80FFFFFF)),
        ),
        elevation: 0,
        actions: [
          if (widget.file.type == 'file' && widget.file.name.endsWith('.txt'))
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: ElevatedButton(
                onPressed: isModified ? saveFile : null,
                style: ElevatedButton.styleFrom(
                  backgroundColor: Color(0x0FFFFFFF),
                ),
                child: Text('Сохранить', style: TextStyle(color: Color(0x80FFFFFF))),
              ),
            ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(0.0), // Убрал padding чтобы текст был вверху
        child: _buildContent(),
      ),
    );
  }

  @override
  void dispose() {
    controller.dispose();
    super.dispose();
  }
}
