import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/core/themes/theme.dart';

import '../../themes/colors.dart';

class CustomButton extends StatelessWidget {
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
  Widget build(BuildContext context) {
    return Semantics(
      label: text ?? 'Кнопка',
      child: GestureDetector(
        onTap: onTap,
        child: Container(
          padding: const EdgeInsets.all(5),
          height: 40,
          decoration: BoxDecoration(
            gradient: gradient ? AppColors.gradient_4 : null,
            color: gradient ? null : AppColors.white15,
            borderRadius: BorderRadius.circular(12),
          ),
          child: Row(
            children: [
              if (prefixIcon != null) Icon(prefixIcon, color: gradient ? AppColors.white : AppColors.white128),
              if (text != null)
                Row(
                  children: [
                    const SizedBox(width: 5),
                    Text(text!, style: gradient ? AppTextStyles.style4 : AppTextStyles.style3),
                  ],
                ),
              if (suffixIcon != null)
                Row(
                  children: [
                    const SizedBox(width: 5),
                    Icon(suffixIcon, color: gradient ? AppColors.white : AppColors.white128),
                  ],
                ),
            ],
          ),
        ),
      ),
    );
  }
}
