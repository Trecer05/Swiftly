import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';

class CustomRadio extends StatelessWidget {
  final bool isChecked;
  final VoidCallback onPressed;  
  const CustomRadio({super.key, required this.isChecked, required this.onPressed});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onPressed,
      child: Container(
        width: 16,
        height: 16,
        decoration: BoxDecoration(
          color: isChecked ? AppColors.blue : AppColors.transparent,
          border: Border.all(
            color: isChecked ? AppColors.blue : AppColors.white,
            width: 2,
          ),
          borderRadius: BorderRadius.circular(8),
        ),
        child: isChecked
            ? const Icon(Icons.check, color: AppColors.white, size: 12)
            : null,
      ),
    );
  }
}
