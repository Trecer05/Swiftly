import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';

class ContainerWidget extends StatelessWidget {
  final String imagePath;
  const ContainerWidget({super.key, required this.imagePath});

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 36,
      padding: EdgeInsets.symmetric(vertical: 10),
      decoration: BoxDecoration(
        color: AppColors.white15,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Center(
        child: SizedBox(
          width: 12,
          height: 12,
          child: Image.asset(imagePath, fit: BoxFit.fill),
        ),
      ),
    );
  }
}