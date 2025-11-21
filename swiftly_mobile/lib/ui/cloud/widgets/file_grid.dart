import 'dart:io';
import 'package:flutter/material.dart';
import 'file_model.dart';

/// Виджет для отображения файлов в адаптивной сетке
class FileGrid extends StatelessWidget {
  final List<FileInfo> files;
  final Function(FileInfo)? onTap;
  final Function(FileInfo)? onDelete;
  final Function(FileInfo, String)? onRename;
  final Function(FileInfo)? onCopy;
  final Function(String)? onCreateFolder;
  final String currentPath;

  const FileGrid({
    required this.files,
    required this.currentPath,
    this.onTap,
    this.onDelete,
    this.onRename,
    this.onCopy,
    this.onCreateFolder,
    super.key,
  });

  /// Вычисляет количество колонок на основе доступной ширины
  int _calculateCrossAxisCount(double width) {
    const itemWidth = 116.0;
    const spacing = 24.0;
    final availableWidth = width - (spacing * 2);
    final count = ((availableWidth + spacing) / (itemWidth + spacing)).floor();
    return count.clamp(1, 20);
  }

  /// Показывает диалог подтверждения удаления
  void _showDeleteConfirmDialog(BuildContext context, FileInfo file) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('Удалить ${file.name}?'),
        content: Text(
          file.type == 'folder'
              ? 'Папка и всё содержимое будут удалены безвозвратно'
              : 'Файл будет удалён безвозвратно',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: Text('Отмена'),
          ),
          TextButton(
            onPressed: () {
              onDelete?.call(file);
              Navigator.pop(ctx);
            },
            child: Text('Удалить', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }

  /// Показывает диалог переименования файла
  void _showRenameDialog(BuildContext context, FileInfo file) {
    final controller = TextEditingController(text: file.name);
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('Переименовать'),
        content: TextField(
          controller: controller,
          decoration: InputDecoration(
            hintText: 'Новое имя',
            border: OutlineInputBorder(),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: Text('Отмена'),
          ),
          TextButton(
            onPressed: () {
              if (controller.text.isNotEmpty) {
                onRename?.call(file, controller.text);
                Navigator.pop(ctx);
              }
            },
            child: Text('Переименовать'),
          ),
        ],
      ),
    );
  }

  /// Показывает диалог создания папки в текущей директории
  void _showCreateFolderInCurrentDialog(BuildContext context) {
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('Создать папку'),
        content: TextField(
          controller: controller,
          decoration: InputDecoration(
            hintText: 'Имя папки',
            border: OutlineInputBorder(),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: Text('Отмена'),
          ),
          TextButton(
            onPressed: () {
              if (controller.text.isNotEmpty) {
                onCreateFolder?.call('$currentPath/${controller.text}');
                Navigator.pop(ctx);
              }
            },
            child: Text('Создать'),
          ),
        ],
      ),
    );
  }

  /// Показывает диалог создания папки внутри родительской папки
  void _showCreateFolderDialog(BuildContext context, FileInfo parentFolder) {
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('Создать папку'),
        content: TextField(
          controller: controller,
          decoration: InputDecoration(
            hintText: 'Имя папки',
            border: OutlineInputBorder(),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: Text('Отмена'),
          ),
          TextButton(
            onPressed: () {
              if (controller.text.isNotEmpty) {
                onCreateFolder?.call('${parentFolder.localPath}/${controller.text}');
                Navigator.pop(ctx);
              }
            },
            child: Text('Создать'),
          ),
        ],
      ),
    );
  }

  /// Показывает контекстное меню для файла
  void _showContextMenu(BuildContext context, TapDownDetails details, FileInfo file) {
    final RenderBox overlay = Overlay.of(context).context.findRenderObject() as RenderBox;

    showMenu<String>(
      context: context,
      position: RelativeRect.fromRect(
        Rect.fromPoints(
          details.globalPosition,
          details.globalPosition.translate(1, 1),
        ),
        Offset.zero & overlay.size,
      ),
      items: [
        PopupMenuItem(
          value: 'rename',
          child: Row(
            children: [
              Icon(Icons.edit, size: 18),
              SizedBox(width: 12),
              Text('Переименовать'),
            ],
          ),
        ),
        PopupMenuItem(
          value: 'copy',
          child: Row(
            children: [
              Icon(Icons.content_copy, size: 18),
              SizedBox(width: 12),
              Text('Копировать'),
            ],
          ),
        ),
        PopupMenuItem(
          value: 'delete',
          child: Row(
            children: [
              Icon(Icons.delete, size: 18, color: Colors.red),
              SizedBox(width: 12),
              Text('Удалить', style: TextStyle(color: Colors.red)),
            ],
          ),
        ),
        if (file.type == 'folder') ...[
          PopupMenuDivider(),
          PopupMenuItem(
            value: 'create_folder',
            child: Row(
              children: [
                Icon(Icons.create_new_folder_outlined, size: 18),
                SizedBox(width: 12),
                Text('Создать папку'),
              ],
            ),
          ),
        ],
      ],
    ).then((value) {
      if (value == null) return;

      switch (value) {
        case 'rename':
          _showRenameDialog(context, file);
          break;
        case 'copy':
          onCopy?.call(file);
          break;
        case 'delete':
          _showDeleteConfirmDialog(context, file);
          break;
        case 'create_folder':
          _showCreateFolderDialog(context, file);
          break;
      }
    });
  }

  /// Показывает контекстное меню для пустой области
  void _showEmptyAreaContextMenu(BuildContext context, TapDownDetails details) {
    final RenderBox overlay = Overlay.of(context).context.findRenderObject() as RenderBox;

    showMenu<String>(
      context: context,
      position: RelativeRect.fromRect(
        Rect.fromPoints(
          details.globalPosition,
          details.globalPosition.translate(1, 1),
        ),
        Offset.zero & overlay.size,
      ),
      items: [
        PopupMenuItem(
          value: 'create_folder',
          child: Row(
            children: [
              Icon(Icons.create_new_folder_outlined, size: 18),
              SizedBox(width: 12),
              Text('Создать папку'),
            ],
          ),
        ),
      ],
    ).then((value) {
      if (value == 'create_folder') {
        _showCreateFolderInCurrentDialog(context);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      behavior: HitTestBehavior.translucent,
      onSecondaryTapDown: (details) => _showEmptyAreaContextMenu(context, details),
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: LayoutBuilder(
          builder: (context, constraints) {
            int crossAxisCount = _calculateCrossAxisCount(constraints.maxWidth);

            return GridView.builder(
              itemCount: files.length,
              gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: crossAxisCount,
                mainAxisSpacing: 0,
                crossAxisSpacing: 24,
                childAspectRatio: 116 / 100,
              ),
              itemBuilder: (_, i) => GestureDetector(
                onTap: () => onTap?.call(files[i]),
                onSecondaryTapDown: (details) => _showContextMenu(context, details, files[i]),
                child: FileGridItem(file: files[i]),
              ),
            );
          },
        ),
      ),
    );
  }
}

/// Виджет для отображения одного файла в сетке
class FileGridItem extends StatelessWidget {
  final FileInfo file;

  const FileGridItem({required this.file, super.key});

  @override
  Widget build(BuildContext context) {
    Widget icon;

    if (file.type == 'image' && file.localPath != null) {
      icon = ClipRRect(
        borderRadius: BorderRadius.circular(8),
        child: Image.file(
          File(file.localPath!),
          fit: BoxFit.cover,
          width: 60,
          height: 60,
        ),
      );
    } else if (file.type == 'folder') {
      icon = Icon(Icons.folder, size: 60, color: Colors.blue[300]);
    } else if (file.type == 'archive') {
      icon = Icon(Icons.folder_zip, size: 60, color: Colors.orange[300]);
    } else {
      icon = Icon(Icons.description, size: 60, color: Colors.grey[400]);
    }

    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        icon,
        SizedBox(height: 8),
        Padding(
          padding: EdgeInsets.symmetric(horizontal: 4),
          child: Text(
            file.name,
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
            textAlign: TextAlign.center,
            style: TextStyle(
              color: Colors.white,
              fontSize: 12,
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
      ],
    );
  }
}
