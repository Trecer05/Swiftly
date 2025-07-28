import 'package:flutter/material.dart';

import '../../../core/themes/colors.dart';
import '../../../core/themes/theme.dart';

class TextFieldWidget extends StatelessWidget {
  final String hintText;
  final TextEditingController controller;

  const TextFieldWidget({
    super.key,
    required this.hintText,
    required this.controller,
  });

  @override
  Widget build(BuildContext context) {
    return TextField(
      controller: controller,
      style: AppTextStyles.style4,
      decoration: InputDecoration(
        filled: true,
        fillColor: AppColors.white15,
        hintText: hintText,
        hintStyle: AppTextStyles.style3,
        enabledBorder: const OutlineInputBorder(
          borderRadius: BorderRadius.all(Radius.circular(12)),
          borderSide: BorderSide.none,
        ),
        focusedBorder: const OutlineInputBorder(
          borderRadius: BorderRadius.all(Radius.circular(12)),
          borderSide: BorderSide(color: Colors.white),
        ),
      ),
    );
  }
}
