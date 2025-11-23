import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter_context_menu/flutter_context_menu.dart';
import 'file_model.dart';
import 'package:desktop_drop/desktop_drop.dart';
import 'package:cross_file/cross_file.dart'; 

class FileGrid extends StatelessWidget {
  final List<FileInfo> files;
  final Function(FileInfo)? onTap;
  final Function(FileInfo)? onDelete;
  final Function(FileInfo, String)? onRename;
  final Function(FileInfo)? onCopy;
  final Function(String)? onCreateFolder;
  final String currentPath;

  // тут настраиваешь сдвиг
  final double menuOffsetX; // >0 — левее, <0 — правее
  final double menuOffsetY; // >0 — выше,  <0 — ниже

  final void Function(List<XFile> files)? onDropFiles;

  const FileGrid({
    required this.files,
    required this.currentPath,
    this.onTap,
    this.onDelete,
    this.onRename,
    this.onCopy,
    this.onCreateFolder,
    this.menuOffsetX = 24,   // можно крутить под себя
    this.menuOffsetY = 12,
    this.onDropFiles,
    super.key,
  });

  int _calculateCrossAxisCount(double width) {
    const itemWidth = 116.0;
    const spacing = 24.0;
    final availableWidth = width - (spacing * 2);
    final count = ((availableWidth + spacing) / (itemWidth + spacing)).floor();
    return count.clamp(1, 20);
  }

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
            child: const Text('Отмена'),
          ),
          TextButton(
            onPressed: () {
              onDelete?.call(file);
              Navigator.pop(ctx);
            },
            child: const Text('Удалить', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }

  void _showRenameDialog(BuildContext context, FileInfo file) {
    final controller = TextEditingController(text: file.name);
    showDialog(
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
            child: const Text('Переименовать'),
          ),
        ],
      ),
    );
  }

  void _showCreateFolderDialog(BuildContext context) {
    final controller = TextEditingController();
    showDialog(
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

  // ПКМ по файлу: показываем flutter_context_menu со сдвигом.[web:7]
  void _showFileMenuWithOffset(
    BuildContext context,
    TapDownDetails details,
    FileInfo file,
  ) {
    final basePos = details.globalPosition;
    final offsetPos = Offset(
      basePos.dx - 390,
      basePos.dy + 10,
    );

    final fileMenu = ContextMenu(
      position: offsetPos, // ключ — задаём вручную позицию меню[web:7]
      entries: <ContextMenuEntry>[
        const MenuHeader(text: 'Файл'),
        MenuItem(
          label: 'Переименовать',
          icon: Icons.edit,
          value: 'rename',
          onSelected: () => _showRenameDialog(context, file),
        ),
        MenuItem(
          label: 'Копировать',
          icon: Icons.content_copy,
          value: 'copy',
          onSelected: () => onCopy?.call(file),
        ),
        const MenuDivider(),
        MenuItem(
          label: 'Удалить',
          icon: Icons.delete,
          value: 'delete',
          onSelected: () => _showDeleteConfirmDialog(context, file),
        ),
      ],
    );

    // способ 1 из доки пакета — показать меню в заданной position[web:7]
    showContextMenu(
      context,
      contextMenu: fileMenu,
    );
  }

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final crossAxisCount = _calculateCrossAxisCount(constraints.maxWidth);

        return DropTarget(
          // сюда прилетают файлы из проводника
          onDragDone: (detail) {
            // detail.files имеет тип List<XFile>
            onDropFiles?.call(detail.files);
          },
          child: GridView.builder(
            padding: const EdgeInsets.fromLTRB(20, 0, 24, 24),
            itemCount: files.length,
            gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: crossAxisCount,
              mainAxisSpacing: 24,
              crossAxisSpacing: 24,
              childAspectRatio: 116 / 100,
            ),
            itemBuilder: (_, i) {
              final file = files[i];

              return GestureDetector(
                onTap: () => onTap?.call(file),
                onSecondaryTapDown: (details) =>
                    _showFileMenuWithOffset(context, details, file),
                child: FileGridItem(file: file),
              );
            },
          ),
        );
      },
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
        const SizedBox(height: 8),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 4),
          child: Text(
            file.name,
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
            textAlign: TextAlign.center,
            style: const TextStyle(
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
