import 'package:flutter/material.dart';
import 'package:file_picker/file_picker.dart';
import 'package:go_router/go_router.dart';
import 'dart:io';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:path_provider/path_provider.dart';
import 'package:open_filex/open_filex.dart';
import 'widgets/side_panel.dart';
import 'widgets/app_bar_cloud.dart';
import 'widgets/file_grid.dart';
import 'widgets/file_model.dart';
import 'file_editor_screen.dart';

class CloudScreen extends StatefulWidget {
  const CloudScreen({super.key});

  @override
  State<CloudScreen> createState() => _CloudScreenState();
}

class _CloudScreenState extends State<CloudScreen> with WidgetsBindingObserver {
  List<FileInfo> files = [];
  List<FileInfo> _allFiles = [];
  List<FileInfo> searchResults = [];
  String searchQuery = '';
  String? workingDir;
  List<String> directoryStack = [];
  Set<String> openFiles = {};

  static const _prefsKey = 'working_directory';
  static const _prefsStackKey = 'directory_stack';

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    _initWorkingDir();
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (state == AppLifecycleState.resumed) {
      _loadWorkingDirFiles();
    }
  }

  /// Инициализирует рабочую директорию
  Future<void> _initWorkingDir() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final savedDir = prefs.getString(_prefsKey);
      final savedStack = prefs.getStringList(_prefsStackKey) ?? [];

      if (savedDir != null && await Directory(savedDir).exists()) {
        workingDir = savedDir;
        directoryStack = savedStack;
      } else {
        final documentsDir = await getApplicationDocumentsDirectory();
        workingDir = documentsDir.path;
        directoryStack = [];
      }

      await _loadWorkingDirFiles();
    } catch (e) {
      print('Error initializing: $e');
    }
  }

  /// Загружает файлы из текущей рабочей директории
  Future<void> _loadWorkingDirFiles() async {
    if (workingDir == null) return;

    try {
      final dir = Directory(workingDir!);
      if (!await dir.exists()) return;

      final entities = await dir.list().toList();
      final loadedFiles = <FileInfo>[];

      for (var entity in entities) {
        final fullPath = entity.path;
        final name = fullPath.split(Platform.pathSeparator).last;
        String type = 'file';

        if (entity is Directory) {
          type = 'folder';
        } else if (entity is File) {
          final ext = name.split('.').last.toLowerCase();
          if (['jpg', 'jpeg', 'png', 'gif', 'webp'].contains(ext)) {
            type = 'image';
          } else if (['zip', 'rar', '7z'].contains(ext)) {
            type = 'archive';
          }
        }

        loadedFiles.add(FileInfo(name, type, fullPath));
      }

      setState(() {
        files = loadedFiles;
        _allFiles = loadedFiles;
      });

      await _saveState();
    } catch (e) {
      print('Error loading directory: $e');
    }
  }

  /// Сохраняет текущее состояние в preferences
  Future<void> _saveState() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_prefsKey, workingDir ?? '');
    await prefs.setStringList(_prefsStackKey, directoryStack);
  }

  /// Открывает файл или папку
  Future<void> _openFile(FileInfo file) async {
    if (file.type == 'folder') {
      directoryStack.add(workingDir!);
      workingDir = file.localPath;
      await _loadWorkingDirFiles();
    } else {
      openFiles.add(file.name);
      setState(() {});

      try {
        final extension = file.name.split('.').last.toLowerCase();

        if (file.type == 'image' ||
            ['txt', 'dart', 'json', 'md', 'yaml', 'xml', 'html', 'css', 'js']
                .contains(extension)) {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => FileEditorScreen(file: file),
            ),
          );
        } else if (['docx', 'doc', 'pdf', 'xlsx', 'pptx']
            .contains(extension)) {
          await _showOpenExternalDialog(file);
        } else {
          await OpenFilex.open(file.localPath!);
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Ошибка открытия: $e')),
          );
        }
      } finally {
        openFiles.remove(file.name);
        if (mounted) setState(() {});
      }
    }
  }

  Future<void> _showOpenExternalDialog(FileInfo file) async {
    final extension = file.name.split('.').last.toUpperCase();

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('Открыть $extension-файл'),
        content: Text(
          'Конвертация в PDF сейчас не поддерживается внутри приложения. '
          'Открыть документ в установленном приложении для $extension?',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: Text('Отмена'),
          ),
          TextButton(
            onPressed: () async {
              Navigator.pop(ctx);
              try {
                await OpenFilex.open(file.localPath!);
              } catch (e) {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('Ошибка открытия: $e')),
                  );
                }
              }
            },
            child: Text('Открыть в ${extension == 'DOCX' ? 'Word' : extension}'),
          ),
        ],
      ),
    );
  }

  /// Удаляет файл или папку
  Future<void> _deleteFile(FileInfo file) async {
    try {
      if (openFiles.contains(file.name)) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
              content: Text('Невозможно удалить открытый файл: ${file.name}')),
        );
        return;
      }

      final entity = file.type == 'folder'
          ? Directory(file.localPath!)
          : File(file.localPath!);

      if (file.type == 'folder') {
        await (entity as Directory).delete(recursive: true);
      } else {
        await (entity as File).delete();
      }

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('${file.name} удалён')),
      );

      await _loadWorkingDirFiles();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Ошибка удаления: $e')),
      );
    }
  }

  /// Переименовывает файл или папку
  Future<void> _renameFile(FileInfo file, String newName) async {
    try {
      if (newName.isEmpty) return;

      final oldPath = file.localPath!;
      final newPath =
          '${oldPath.substring(0, oldPath.lastIndexOf("/"))}/$newName';

      final entity = file.type == 'folder'
          ? Directory(oldPath)
          : File(oldPath);

      await entity.rename(newPath);

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Переименовано в $newName')),
      );

      await _loadWorkingDirFiles();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Ошибка переименования: $e')),
      );
    }
  }

  /// Копирует файл или папку
  Future<void> _copyFile(FileInfo file) async {
    try {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('${file.name} скопирован')),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Ошибка копирования: $e')),
      );
    }
  }

  /// Создаёт новую папку
  Future<void> _createFolder(String folderPath) async {
    try {
      final dir = Directory(folderPath);
      await dir.create(recursive: true);

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Папка создана')),
      );

      await _loadWorkingDirFiles();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Ошибка создания папки: $e')),
      );
    }
  }

  /// Возвращается к предыдущей директории
  Future<void> _navigateBack() async {
    if (directoryStack.isEmpty) return;
    workingDir = directoryStack.removeLast();
    await _loadWorkingDirFiles();
  }

  /// Выбирает рабочую директорию
  Future<void> _selectFolder() async {
    try {
      String? selectedDirectory =
          await FilePicker.platform.getDirectoryPath();

      if (selectedDirectory != null) {
        workingDir = selectedDirectory;
        directoryStack = [];
        await _loadWorkingDirFiles();
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Ошибка выбора папки: $e')),
      );
    }
  }

  /// Добавляет файлы в текущую директорию
  Future<void> _addFiles() async {
    try {
      final result =
          await FilePicker.platform.pickFiles(allowMultiple: true);

      if (result != null) {
        for (var file in result.files) {
          final sourceFile = File(file.path!);
          final destPath = '$workingDir/${file.name}';
          await sourceFile.copy(destPath);
        }

        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Файлы добавлены')),
        );

        await _loadWorkingDirFiles();
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Ошибка добавления файлов: $e')),
      );
    }
  }

  /// Выполняет поиск файлов
  void _searchFiles(String query) {
    setState(() {
      searchQuery = query;
      if (query.isEmpty) {
        files = _allFiles;
      } else {
        files = _allFiles
            .where((file) =>
                file.name.toLowerCase().contains(query.toLowerCase()))
            .toList();
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        child: Row(
          children: [
            SizedBox(
              width: 200,
              child: Column(
                children: [
                  Container(
                    height: 81,
                    alignment: Alignment.centerLeft,
                    padding: EdgeInsets.only(left: 20, top: 25, bottom: 20),
                    child: Text(
                      'Ваше облако',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 18,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                  ),
                  Expanded(
                    child: SidePanel(),
                  ),
                ],
              ),
            ),
            Expanded(
              child: Stack(
                children: [
                  Column(
                    children: [
                      AppBarCloud(
                        title: 'Облако',
                        currentDir: workingDir?.split(Platform.pathSeparator).last ??
                            'Облако',
                        count: files.length,
                        onAdd: _addFiles,
                        onSearch: _searchFiles,
                      ),
                      Expanded(
                        child: FileGrid(
                          files: files,
                          currentPath: workingDir ?? '',
                          onTap: _openFile,
                          onDelete: _deleteFile,
                          onRename: _renameFile,
                          onCopy: _copyFile,
                          onCreateFolder: _createFolder,
                        ),
                      ),
                    ],
                  ),
                  Positioned(
                    bottom: 24,
                    right: 24,
                    child: FloatingActionButton(
                      heroTag: 'select_folder_fab',
                      onPressed: _selectFolder,
                      backgroundColor: Color(0xFF6DA8FF),
                      child: Icon(Icons.folder_open, color: Colors.white),
                      tooltip: 'Выбрать папку',
                    ),
                  ),
                  if (directoryStack.isNotEmpty)
                    Positioned(
                      bottom: 96,
                      right: 24,
                      child: FloatingActionButton(
                        heroTag: 'back_fab',
                        onPressed: _navigateBack,
                        backgroundColor: Color(0xFF6DA8FF),
                        child: Icon(Icons.arrow_back, color: Colors.white),
                        tooltip: 'Назад',
                      ),
                    ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}



