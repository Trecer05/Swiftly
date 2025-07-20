import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/core/themes/theme.dart';

import '../../core/themes/colors.dart';

class CustomButton extends StatelessWidget {
  final IconData? prefixIcon;
  final String? text;
  final VoidCallback onTap;
  const CustomButton({
    super.key,
    this.prefixIcon,
    this.text,
    required this.onTap,
  }) : assert(text != null || prefixIcon != null);

  @override
  Widget build(BuildContext context) {
    return Semantics(
      label: text ?? 'Кнопка',
      child: GestureDetector(
        onTap: onTap,
        child: Container(
          padding: const EdgeInsets.all(5),
          height: 40,
          decoration: BoxDecoration(
            color: AppColors.white15,
            borderRadius: BorderRadius.circular(12),
          ),
          child: Row(
            children: [
              if (prefixIcon != null) Icon(prefixIcon, color: AppColors.white128),
              if (text != null)
                Row(
                  children: [
                    const SizedBox(width: 5),
                    Text(text!, style: AppTextStyles.style3),
                  ],
                ),
            ],
          ),
        ),
      ),
    );
  }
}
