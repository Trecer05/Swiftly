import 'dart:io';

import 'package:flutter/material.dart';

import 'widgets/file_model.dart';      // FileInfo
import 'widgets/app_bar_cloud.dart';   // AppBarCloud
import 'widgets/file_grid.dart';       // FileGrid
import 'widgets/side_panel.dart';      // SidePanel

class CloudDesktopScreen extends StatelessWidget {
  final List<FileInfo> files;
  final String? workingDir;
  final List directoryStack;

  final Function(FileInfo) onTap;
  final Function(FileInfo) onDelete;
  final Function(FileInfo, String) onRename;
  final Function(FileInfo) onCopy;
  final Function(String) onCreateFolder;
  final VoidCallback onAddFiles;
  final Function(String) onSearchFiles;
  final VoidCallback onSelectFolder;
  final VoidCallback onNavigateBack;

  const CloudDesktopScreen({
    super.key,
    required this.files,
    required this.workingDir,
    required this.directoryStack,
    required this.onTap,
    required this.onDelete,
    required this.onRename,
    required this.onCopy,
    required this.onCreateFolder,
    required this.onAddFiles,
    required this.onSearchFiles,
    required this.onSelectFolder,
    required this.onNavigateBack,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Row(
        children: [
          SizedBox(
            width: 200,
            child: Column(
              children: [
                Container(
                  height: 81,
                  alignment: Alignment.centerLeft,
                  padding: const EdgeInsets.only(
                    left: 20,
                    top: 25,
                    bottom: 20,
                  ),
                  child: const Text(
                    'Ваше облако',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 18,
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                ),
                const Expanded(
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
                      currentDir: workingDir?.split(
                            Platform.pathSeparator,
                          ).last ??
                          'Облако',
                      count: files.length,
                      onAdd: onAddFiles,
                      onSearch: onSearchFiles,
                    ),
                    Expanded(
                      child: FileGrid(
                        files: files,
                        currentPath: workingDir ?? '',
                        onTap: onTap,
                        onDelete: onDelete,
                        onRename: onRename,
                        onCopy: onCopy,
                        onCreateFolder: onCreateFolder,
                      ),
                    ),
                  ],
                ),
                Positioned(
                  bottom: 24,
                  right: 24,
                  child: FloatingActionButton(
                    heroTag: 'select_folder_fab',
                    onPressed: onSelectFolder,
                    backgroundColor: const Color(0xFF6DA8FF),
                    child: const Icon(
                      Icons.folder_open,
                      color: Colors.white,
                    ),
                    tooltip: 'Выбрать папку',
                  ),
                ),
                if (directoryStack.isNotEmpty)
                  Positioned(
                    bottom: 96,
                    right: 24,
                    child: FloatingActionButton(
                      heroTag: 'back_fab',
                      onPressed: onNavigateBack,
                      backgroundColor: const Color(0xFF6DA8FF),
                      child: const Icon(
                        Icons.arrow_back,
                        color: Colors.white,
                      ),
                      tooltip: 'Назад',
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
