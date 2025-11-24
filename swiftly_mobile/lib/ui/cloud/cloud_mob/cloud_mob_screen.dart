import 'package:flutter/material.dart';
import '../widgets/file_model.dart'; // FileInfo
import '../widgets/app_bar_cloud.dart'; // для CloudSearchField / CloudAddButton
import '../widgets/file_grid.dart';    // для FileGridItem (иконка + подпись)
import '../cloud_screen.dart';

class CloudTabMobile extends StatelessWidget {
  final List<FileInfo> files;
  final String currentPath;

  final void Function(FileInfo)? onTap;
  final void Function(FileInfo)? onDelete;
  final void Function(FileInfo, String)? onRename;
  final void Function(FileInfo)? onShare;
  final void Function(String)? onCreateFolder;
  final VoidCallback? onAddFile;
  final void Function(String)? onSearch;

  const CloudTabMobile({
    required this.files,
    required this.currentPath,
    this.onTap,
    this.onDelete,
    this.onRename,
    this.onShare,
    this.onAddFile,
    this.onSearch,
    this.onCreateFolder,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF050816),
      body: SafeArea(
        child: Column(
          children: [
            // верхняя панель как в макете
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 8, 16, 8),
              child: Row(
                children: [
                  Text(
                    'Файлы',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                  const SizedBox(width: 4),
                  Text(
                    '(${files.length})',
                    style: const TextStyle(
                      color: Color(0xFF6DA8FF),
                      fontSize: 11,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  const Spacer(),
                  // переиспользуем поиск и кнопку добавления из app_bar_cloud.dart
                  CloudSearchField(onSearch: onSearch),
                  const SizedBox(width: 8),
                  CloudAddButton(onAdd: onAddFile ?? () {}),
                ],
              ),
            ),
            const SizedBox(height: 8),
            Expanded(
              child: _MobileFileGrid(
                files: files,
                currentPath: currentPath,
                onTap: onTap,
                onDelete: onDelete,
                onRename: onRename,
                onCreateFolder: onCreateFolder,
                onShare: onShare,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
class _MobileFileGrid extends StatelessWidget {
  final List<FileInfo> files;
  final String currentPath;

  final void Function(FileInfo)? onTap;
  final void Function(FileInfo)? onDelete;
  final void Function(FileInfo, String)? onRename;
  final void Function(String)? onCreateFolder;
  final void Function(FileInfo)? onShare;

  const _MobileFileGrid({
    required this.files,
    required this.currentPath,
    this.onTap,
    this.onDelete,
    this.onRename,
    this.onCreateFolder,
    this.onShare,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return GridView.builder(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 3,
        crossAxisSpacing: 12,
        mainAxisSpacing: 12,
        childAspectRatio: 0.75,
      ),
      itemCount: files.length,
      itemBuilder: (context, index) {
        final file = files[index];
        return GestureDetector(
          onTap: () => onTap?.call(file),
          onLongPress: () => _showFileActions(context, file),
          child: FileGridItem(file: file),
        );
      },

      
    );
  }

  void _showFileActions(BuildContext context, FileInfo file) {
    showModalBottomSheet(
      context: context,
      backgroundColor: const Color(0xFF111827),
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (ctx) {
        return SafeArea(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              ListTile(
                leading: const Icon(Icons.drive_file_rename_outline,
                    color: Colors.white),
                title: const Text('Переименовать',
                    style: TextStyle(color: Colors.white)),
                onTap: () async {
                  Navigator.pop(ctx);
                  await _showRenameDialog(context, file);
                },
              ),
              ListTile(
                leading:
                    const Icon(Icons.copy, color: Colors.white),
                title: const Text('Поделиться',
                    style: TextStyle(color: Colors.white)),
                onTap: () {
                  Navigator.pop(ctx);
                  onShare?.call(file);   // ← в буфер (логика в _copyFile)
                },
              ),
              ListTile(
                leading: const Icon(Icons.delete, color: Colors.red),
                title: const Text('Удалить',
                    style: TextStyle(color: Colors.red)),
                onTap: () {
                  Navigator.pop(ctx);
                  onDelete?.call(file);
                },
              ),
            ],
          ),
        );
      },
    );
  }

  Future<void> _showRenameDialog(
      BuildContext context, FileInfo file) async {
    final controller = TextEditingController(text: file.name);
    await showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Переименовать'),
        content: TextField(
          controller: controller,
          decoration: const InputDecoration(
            hintText: 'Новое имя',
            border: OutlineInputBorder(),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Отмена'),
          ),
          TextButton(
            onPressed: () {
              if (controller.text.isNotEmpty) {
                onRename?.call(file, controller.text);
                Navigator.pop(ctx);
              }
            },
            child: const Text('Сохранить'),
          ),
        ],
      ),
    );
  }

  Future<void> _showCreateFolderDialog(BuildContext context) async {
    final controller = TextEditingController();
    await showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Создать папку'),
        content: TextField(
          controller: controller,
          decoration: const InputDecoration(
            hintText: 'Имя папки',
            border: OutlineInputBorder(),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Отмена'),
          ),
          TextButton(
            onPressed: () {
              if (controller.text.isNotEmpty) {
                onCreateFolder?.call('$currentPath/${controller.text}');
                Navigator.pop(ctx);
              }
            },
            child: const Text('Создать'),
          ),
        ],
      ),
    );
  }
}
