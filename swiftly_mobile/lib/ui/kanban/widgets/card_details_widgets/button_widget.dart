import 'package:flutter/material.dart';

import '../../../core/themes/colors.dart';

class ButtonWidget extends StatefulWidget {
  final String title;
  final Color color;
  final VoidCallback onTap;
  const ButtonWidget({super.key, required this.title, required this.color, required this.onTap});

  @override
  State<ButtonWidget> createState() => _ButtonWidgetState();
}

class _ButtonWidgetState extends State<ButtonWidget> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => isHovered = true),
      onExit: (_) => setState(() => isHovered = false),
      child: GestureDetector(
        onTap: widget.onTap,
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
          decoration: BoxDecoration(
            color: widget.color,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: isHovered ? AppColors.white : AppColors.transparent,
            ),
          ),
          child: Text(widget.title),
        ),
      ),
    );
  }
}
