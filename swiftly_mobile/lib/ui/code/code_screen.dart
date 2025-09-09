import 'package:flutter/material.dart';
import 'dart:ui';
import 'dart:io';
import 'package:file_picker/file_picker.dart';

final ValueNotifier<String> activeTab = ValueNotifier("CameraViewModel.swift");
final ValueNotifier<List<String>> openTabs = ValueNotifier(["CameraViewModel.swift"]);
final ValueNotifier<Directory?> rootDir = ValueNotifier<Directory?>(null);
final ValueNotifier<String?> activeFilePath = ValueNotifier<String?>(null);
final Map<String, String> openedFiles = {}; // path -> content
final Map<String, TextEditingController> editors = {}; // path -> controller

final ValueNotifier<double> explorerWidth = ValueNotifier<double>(260);

class _Acrylic extends StatelessWidget {
  final Widget child;
  final double blur;
  final Color tint;
  final double opacity;
  final BorderRadius? borderRadius;
  const _Acrylic({
    required this.child,
    this.blur = 24,
    this.tint = const Color(0xFF0F172A),
    this.opacity = 0.6,
    this.borderRadius,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return ClipRRect(
      borderRadius: borderRadius ?? BorderRadius.zero,
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: blur, sigmaY: blur),
        child: ColoredBox(
          color: tint.withOpacity(opacity),
          child: child,
        ),
      ),
    );
  }
}

class CodeScreen extends StatelessWidget {
  const CodeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0F172A),
      body: Row(
        children: [
          ValueListenableBuilder<double>(
            valueListenable: explorerWidth,
            builder: (context, w, _) {
              return SizedBox(
                width: w,
                child: const _ExplorerPanel(),
              );
            },
          ),
          MouseRegion(
            cursor: SystemMouseCursors.resizeLeftRight,
            child: GestureDetector(
              behavior: HitTestBehavior.translucent,
              onHorizontalDragUpdate: (d) {
                final next = (explorerWidth.value + d.delta.dx).clamp(200.0, 520.0);
                explorerWidth.value = next;
              },
              child: Container(
                width: 6,
                color: const Color(0x0011182a),
                child: const VerticalDivider(width: 6, color: Color(0xFF334155)),
              ),
            ),
          ),
          const _EditorSide(),
        ],
      ),
    );
  }
}

class _ExplorerPanel extends StatelessWidget {
  const _ExplorerPanel();

  @override
  Widget build(BuildContext context) {
    return _Acrylic(
      tint: const Color(0xFF0F172A),
      opacity: 0.55,
      child: Container(
        decoration: const BoxDecoration(
          color: Colors.transparent,
          border: Border(
            right: BorderSide(color: Color(0xFF334155), width: 1),
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              padding: const EdgeInsets.fromLTRB(16, 14, 16, 12),
              child: const Text(
                'Проводник',
                style: TextStyle(
                  fontSize: 13,
                  color: Color(0xFF94A3B8),
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
            const Divider(height: 1, color: Color(0xFF334155)),
            Expanded(
              child: ValueListenableBuilder<Directory?>(
                valueListenable: rootDir,
                builder: (context, dir, _) {
                  if (dir == null) {
                    return Center(
                      child: Text(
                        'Откройте папку, чтобы увидеть файлы',
                        style: const TextStyle(color: Color(0xFF94A3B8)),
                      ),
                    );
                  }
                  return _DirectoryNode(directory: dir, initiallyExpanded: true);
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _DirectoryNode extends StatefulWidget {
  final Directory directory;
  final bool initiallyExpanded;
  const _DirectoryNode({required this.directory, this.initiallyExpanded = false, super.key});
  @override
  State<_DirectoryNode> createState() => _DirectoryNodeState();
}

class _DirectoryNodeState extends State<_DirectoryNode> {
  late bool _expanded;
  late String _name;
  List<FileSystemEntity> _children = [];

  @override
  void initState() {
    super.initState();
    _expanded = widget.initiallyExpanded;
    _name = widget.directory.path.split(Platform.pathSeparator).where((e) => e.isNotEmpty).last;
    _safeLoad();
  }

  void _safeLoad() {
    try {
      _children = widget.directory
          .listSync()
          .where((e) => !e.path.split(Platform.pathSeparator).last.startsWith('.'))
          .toList()
        ..sort((a, b) {
          final ad = a is Directory;
          final bd = b is Directory;
          if (ad != bd) return ad ? -1 : 1;
          return a.path.toLowerCase().compareTo(b.path.toLowerCase());
        });
    } catch (_) {
      _children = [];
    }
  }

  @override
  Widget build(BuildContext context) {
    final textStyle = const TextStyle(color: Color(0xFFE2E8F0), fontSize: 13, height: 1.3);
    return Theme(
      data: Theme.of(context).copyWith(dividerColor: Colors.transparent, splashColor: const Color.fromRGBO(51, 65, 85, 0.25), hoverColor: const Color.fromRGBO(51, 65, 85, 0.20)),
      child: ListTileTheme(
        dense: true,
        minVerticalPadding: 0,
        horizontalTitleGap: 8,
        child: ExpansionTile(
          tilePadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 0),
          initiallyExpanded: _expanded,
          onExpansionChanged: (v) => setState(() => _expanded = v),
          trailing: const SizedBox.shrink(),
          leading: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              AnimatedRotation(
                turns: _expanded ? 0.25 : 0.0,
                duration: const Duration(milliseconds: 200),
                child: const Icon(Icons.arrow_right, size: 18, color: Color(0xFF94A3B8)),
              ),
              const SizedBox(width: 2),
              const Icon(Icons.folder, size: 16, color: Color(0xFF9CA3AF)),
            ],
          ),
          iconColor: const Color(0xFF94A3B8),
          collapsedIconColor: const Color(0xFF94A3B8),
          title: Text(_name, style: textStyle),
          childrenPadding: const EdgeInsets.only(left: 20),
          children: [
            for (final e in _children)
              if (e is Directory)
                _DirectoryNode(directory: e)
              else
                _FSFileNode(path: e.path),
          ],
        ),
      ),
    );
  }
}

class _FSFileNode extends StatelessWidget {
  final String path;
  const _FSFileNode({required this.path, super.key});

  String get name => path.split(Platform.pathSeparator).last;

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<String?>(
      valueListenable: activeFilePath,
      builder: (context, current, _) {
        final isActive = current == path;
        return InkWell(
          onTap: () async {
            try {
              final content = await File(path).readAsString();
              openedFiles[path] = content;
              editors[path]?.dispose();
              editors[path] = TextEditingController(text: content);
              activeFilePath.value = path;
            } catch (e) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text('Не удалось открыть: ${name}\n${e.runtimeType}')),
              );
            }
          },
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
            decoration: BoxDecoration(color: isActive ? const Color.fromRGBO(37, 58, 102, 0.25) : null),
            child: Row(
              children: [
                const Icon(Icons.description_outlined, size: 16, color: Color(0xFF9CA3AF)),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    name,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(color: Color(0xFFE5E7EB), fontSize: 13),
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}

class _EditorSide extends StatelessWidget {
  const _EditorSide();

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: Column(
        children: const [
          _TopCommandBar(),
          Expanded(child: _EditorArea()),
        ],
      ),
    );
  }
}

class _TopCommandBar extends StatelessWidget {
  const _TopCommandBar();

  Widget _vSep([double h = 18]) => Container(
        width: 1,
        height: h,
        margin: const EdgeInsets.symmetric(horizontal: 12),
        color: const Color(0xFF334155),
      );

  @override
  Widget build(BuildContext context) {
    return _Acrylic(
      tint: const Color(0xFF0F172A),
      opacity: 0.5,
      child: Container(
        color: Colors.transparent,
        padding: const EdgeInsets.fromLTRB(16, 10, 16, 12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
          Align(
            alignment: Alignment.centerLeft,
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                decoration: BoxDecoration(
                  color: const Color.fromRGBO(30, 41, 59, 0.6),
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const _Crumb(text: 'Project'),
                    _vSep(16),
                    const _Crumb(text: 'Name', isActive: true),
                    _vSep(16),
                    const _Crumb(text: 'Working'),
                  ],
                ),
              ),
            ),
          ),
          const SizedBox(height: 10),
          Align(
            alignment: Alignment.centerLeft,
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Row(
                children: [
                  _PillButton(icon: Icons.folder_open, label: 'Open Folder', onPressed: () async {
                    final path = await FilePicker.platform.getDirectoryPath();
                    if (path != null) rootDir.value = Directory(path);
                  }),
                  const SizedBox(width: 10),
                  _PillButton(icon: Icons.add, label: 'Code', onPressed: () {}),
                  const SizedBox(width: 10),
                  _PillButton(icon: Icons.play_arrow, label: 'Run All', onPressed: () {}),
                  const SizedBox(width: 10),
                  _PillButton(icon: Icons.refresh, label: 'Restart', onPressed: () {}),
                  const SizedBox(width: 10),
                  _PillButton(icon: Icons.bug_report_outlined, label: 'Clear all outputs', onPressed: () {}),
                  const SizedBox(width: 10),
                  _PillButton(icon: Icons.more_horiz, label: 'Variables', onPressed: () {}),
                  const SizedBox(width: 10),
                  _PillButton(icon: Icons.more_horiz, label: 'Outline', onPressed: () {}),
                  const SizedBox(width: 10),
                  _PillButton(icon: Icons.more_horiz, label: 'Terminal', onPressed: () {}),
                  const SizedBox(width: 10),
                  _PillButton(icon: Icons.save_outlined, label: 'Save', onPressed: () async {
                    final path = activeFilePath.value;
                    if (path != null && editors[path] != null) {
                      await File(path).writeAsString(editors[path]!.text);
                    }
                  }),
                  const SizedBox(width: 10),
                ],
              ),
            ),
          ),
        ],
      ),
    ),
    );
  }
}

class _Crumb extends StatelessWidget {
  final String text;
  final bool isActive;
  const _Crumb({required this.text, this.isActive = false});
  @override
  Widget build(BuildContext context) {
    if (!isActive) {
      return Text(
        text,
        style: const TextStyle(
          color: Color(0xFFE2E8F0),
          fontSize: 12,
          fontWeight: FontWeight.w600,
        ),
      );
    }
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
      decoration: BoxDecoration(
        color: const Color.fromRGBO(51, 65, 85, 0.55),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        text,
        style: const TextStyle(
          color: Color(0xFFE2E8F0),
          fontSize: 12,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }
}

class _PillButton extends StatelessWidget {
  final String label;
  final IconData icon;
  final VoidCallback onPressed;
  const _PillButton({required this.icon, required this.label, required this.onPressed});

  @override
  Widget build(BuildContext context) {
    return InkWell(
      borderRadius: BorderRadius.circular(12),
      onTap: onPressed,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
        decoration: BoxDecoration(
          color: const Color.fromRGBO(30, 41, 59, 0.6),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Row(
          children: [
            Icon(icon, size: 18, color: const Color(0xFFE2E8F0)),
            const SizedBox(width: 8),
            Text(
              label,
              style: const TextStyle(
                color: Color(0xFFE2E8F0),
                fontSize: 13,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _EditorArea extends StatefulWidget {
  const _EditorArea();
  @override
  State<_EditorArea> createState() => _EditorAreaState();
}

class _EditorAreaState extends State<_EditorArea> {
  final ScrollController _editorScroll = ScrollController();
  final ScrollController _hScroll = ScrollController();

  @override
  void dispose() {
    _editorScroll.dispose();
    _hScroll.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return _Acrylic(
      tint: const Color.fromRGBO(2, 10, 23, 1),
      opacity: 0.45,
      child: Container(
        color: Colors.transparent,
        child: ValueListenableBuilder<String?>(
          valueListenable: activeFilePath,
          builder: (context, path, _) {
          final controller = (path != null)
              ? (editors[path] ?? (editors[path] = TextEditingController(text: openedFiles[path] ?? '')))
              : (editors['__empty'] ?? (editors['__empty'] = TextEditingController(text: '// Откройте файл')));

          return AnimatedBuilder(
            animation: controller,
            builder: (context, __) {
              final int linesCount = controller.text.split('\n').length + 1; // extra caret line
              final lines = List.generate(linesCount, (index) => '${index + 1}'.padLeft(3));
              final double _lineFontSize = 14.0; // same as editor
              final double _lineHeight = 1.5;    // same as editor
              final double _linePixelHeight = _lineFontSize * _lineHeight;
              return LayoutBuilder(
                builder: (context, constraints) {
                  final double _lineFontSize = 14.0; // must match the editor style
                  final double _lineHeight = 1.5;    // must match the editor style
                  final double _linePixelHeight = _lineFontSize * _lineHeight;
                  final int minVisibleLines = (constraints.maxHeight / _linePixelHeight).ceil();
                  final int renderLines = lines.length < minVisibleLines ? minVisibleLines : lines.length;

                  final List<String> _allLines = controller.text.split('\n');
                  final int _maxCols = _allLines.isEmpty ? 1 : _allLines.map((e) => e.length).reduce((a, b) => a > b ? a : b);
                  final double _charWidth = _linePixelHeight * (1.0 / _lineHeight) * 0.62;
                  final double _minEditorWidth = constraints.maxWidth - 48 - 16;
                  final double _contentWidth = (_maxCols + 2) * _charWidth;
                  final double _editorWidth = _contentWidth > _minEditorWidth ? _contentWidth : _minEditorWidth;

                  return SingleChildScrollView(
                    controller: _editorScroll,
                    padding: const EdgeInsets.fromLTRB(0, 0, 16, 16),
                    child: ConstrainedBox(
                      constraints: BoxConstraints(minHeight: constraints.maxHeight),
                      child: Row(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Container(
                            width: 48,
                            decoration: const BoxDecoration(
                              color: Color.fromRGBO(37, 58, 102, 0.25),
                              border: Border(
                                right: BorderSide(color: Color(0xFF334155), width: 1),
                              ),
                            ),
                            padding: const EdgeInsets.symmetric(horizontal: 8),
                            child: ConstrainedBox(
                              constraints: BoxConstraints(minHeight: constraints.maxHeight),
                              child: Align(
                                alignment: Alignment.topRight,
                                child: Text(
                                  List.generate(renderLines, (i) => '${i + 1}').join('\n'),
                                  textAlign: TextAlign.right,
                                  softWrap: false,
                                  style: const TextStyle(
                                    color: Colors.white38,
                                    fontSize: 14,
                                    height: 1.5,
                                    fontFamily: 'SourceCodePro',
                                  ),
                                ),
                              ),
                            ),
                          ),
                          Expanded(
                            child: SingleChildScrollView(
                              scrollDirection: Axis.horizontal,
                              controller: _hScroll,
                              child: SizedBox(
                                width: _editorWidth,
                                child: TextField(
                                  controller: controller,
                                  scrollPhysics: const NeverScrollableScrollPhysics(),
                                  maxLines: null,
                                  keyboardType: TextInputType.multiline,
                                  style: const TextStyle(
                                    fontFamily: 'SourceCodePro',
                                    fontSize: 14,
                                    color: Color(0xFFF1F5F9),
                                    height: 1.5,
                                  ),
                                  decoration: const InputDecoration(border: InputBorder.none, isCollapsed: true),
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  );
                },
              );
            },
          );
          },
        ),
      ),
    );
  }
}
