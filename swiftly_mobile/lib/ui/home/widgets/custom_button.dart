import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/core/themes/theme.dart';

import '../../core/themes/colors.dart';

class CustomButton extends StatelessWidget {
  const CustomButton({super.key});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {},
      child: Container(
        padding: const EdgeInsets.all(5),
        decoration: BoxDecoration(
          color: AppColors.white15,
          borderRadius: BorderRadius.circular(12),
        ),
        child: const Row(
          children: [
            Icon(Icons.add, color: AppColors.white128),
            SizedBox(width: 5),
            Text('Пригласить', style: AppTextStyles.style3),
          ],
        ),
      ),
    );
  }
}