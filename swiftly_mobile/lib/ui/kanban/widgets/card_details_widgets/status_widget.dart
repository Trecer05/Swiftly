import 'package:flutter/material.dart';
import '../../../core/themes/colors.dart';
import '../kanban_status.dart';

class StatusWidget extends StatefulWidget {
  final String statusId;
  final ValueChanged<String> onChanged;

  const StatusWidget({
    super.key,
    required this.statusId,
    required this.onChanged,
  });

  @override
  State<StatusWidget> createState() => _StatusWidgetState();
}

class _StatusWidgetState extends State<StatusWidget> {
  late String _selectedStatusId;
  final LayerLink _layerLink = LayerLink();
  OverlayEntry? _overlayEntry;

  @override
  void initState() {
    super.initState();
    _selectedStatusId = widget.statusId;
  }

  void _toggleDropdown() {
    if (_overlayEntry == null) {
      _showOverlay();
    } else {
      _removeOverlay();
    }
  }

  void _showOverlay() {
    final overlay = Overlay.of(context);
    _overlayEntry = OverlayEntry(
      builder: (context) {
        return Positioned(
          width: 150,
          child: CompositedTransformFollower(
            link: _layerLink,
            showWhenUnlinked: false,
            offset: const Offset(0, 40),
            child: Material(
              color: Colors.transparent,
              child: Container(
                width: null,
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: AppColors.white15,
                  borderRadius: BorderRadius.circular(5),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: List.generate(KanbanStatus.values.length, (index) {
                    final col = KanbanStatus.values[index];
                    final isLast = index == KanbanStatus.values.length - 1;

                    return Column(
                      children: [
                        StatusItemWidget(
                          title: col.title,
                          isSelected: col.id == _selectedStatusId,
                          onTap: () {
                            setState(() {
                              _selectedStatusId = col.id;
                            });
                            widget.onChanged(col.id);
                            _removeOverlay();
                          },
                        ),
                        if (!isLast) const SizedBox(height: 5),
                      ],
                    );
                  }),
                ),
              ),
            ),
          ),
        );
      },
    );

    overlay.insert(_overlayEntry!);
  }

  void _removeOverlay() {
    _overlayEntry?.remove();
    _overlayEntry = null;
  }

  @override
  void dispose() {
    _removeOverlay();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final currentTitle =
        KanbanStatus.values
            .firstWhere((col) => col.id == _selectedStatusId)
            .title;

    return CompositedTransformTarget(
      link: _layerLink,
      child: GestureDetector(
        onTap: _toggleDropdown,
        child: Container(
          height: 36,
          padding: const EdgeInsets.symmetric(horizontal: 10),
          decoration: const BoxDecoration(color: AppColors.white15),
          alignment: Alignment.centerLeft,
          child: Text(
            currentTitle,
            style: const TextStyle(color: AppColors.white),
          ),
        ),
      ),
    );
  }
}

class StatusItemWidget extends StatefulWidget {
  final String title;
  final bool isSelected;
  final VoidCallback onTap;

  const StatusItemWidget({
    super.key,
    required this.title,
    required this.isSelected,
    required this.onTap,
  });

  @override
  State<StatusItemWidget> createState() => _StatusItemWidgetState();
}

class _StatusItemWidgetState extends State<StatusItemWidget> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => isHovered = true),
      onExit: (_) => setState(() => isHovered = false),
      child: GestureDetector(
        onTap: widget.onTap,
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
          decoration: BoxDecoration(
            color: widget.isSelected ? Colors.white24 : Colors.transparent,
            borderRadius: BorderRadius.circular(4),
            border: Border.all(
              color: isHovered ? AppColors.white : AppColors.transparent,
            ),
          ),
          child: Text(
            widget.title,
            style: TextStyle(
              color: widget.isSelected ? Colors.white : Colors.blue,
            ),
          ),
        ),
      ),
    );
  }
}
