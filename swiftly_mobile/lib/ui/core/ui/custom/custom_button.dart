import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/core/themes/theme.dart';

import '../../themes/colors.dart';

class CustomButton extends StatefulWidget {
  final IconData? prefixIcon;
  final String? text;
  final IconData? suffixIcon;
  final bool gradient;
  final VoidCallback onTap;
  const CustomButton({
    super.key,
    this.prefixIcon,
    this.text,
    this.suffixIcon,
    required this.gradient,
    required this.onTap,
  }) : assert(text != null || prefixIcon != null || suffixIcon != null);

  @override
  State<CustomButton> createState() => _CustomButtonState();
}

class _CustomButtonState extends State<CustomButton> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => isHovered = true),
      onExit: (_) => setState(() => isHovered = false),
      child: Semantics(
        label: widget.text ?? 'Кнопка',
        child: GestureDetector(
          onTap: widget.onTap,
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            height: 40,
            decoration: BoxDecoration(
              gradient: widget.gradient ? AppColors.gradient_4 : null,
              color: widget.gradient ? null : isHovered ? AppColors.white38 : AppColors.white15,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(
                color: isHovered ? AppColors.white : AppColors.transparent,
              ),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                if (widget.prefixIcon != null)
                  Icon(
                    widget.prefixIcon,
                    size: 24,
                    color:
                        widget.gradient ? AppColors.white : isHovered ? AppColors.white : AppColors.white128,
                  ),
                if (widget.text != null)
                  Row(
                    children: [
                      const SizedBox(width: 5),
                      Text(
                        widget.text!,
                        style:
                            widget.gradient
                                ? AppTextStyles.style17
                                : AppTextStyles.style7,
                      ),
                    ],
                  ),
                if (widget.suffixIcon != null)
                  Row(
                    children: [
                      const SizedBox(width: 5),
                      Icon(
                        widget.suffixIcon,
                        size: 24,
                        color:
                            widget.gradient
                                ? AppColors.white
                                : AppColors.white128,
                      ),
                    ],
                  ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
