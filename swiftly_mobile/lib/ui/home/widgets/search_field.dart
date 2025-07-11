import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';

class SearchField extends StatelessWidget {
  final String hintText;
  const SearchField({
    super.key,
    required this.hintText,
  });

  @override
  Widget build(BuildContext context) {
    return TextField(
      style: AppTextStyles.text_4,
      decoration: InputDecoration(
        filled: true,
        fillColor: AppColors.white15,
        prefixIcon: const Icon(Icons.search, color: AppColors.white128),
        hintText: hintText,
        hintStyle: AppTextStyles.text_3,

        enabledBorder: const OutlineInputBorder(
          borderRadius: BorderRadius.all(Radius.circular(12)),
          borderSide: BorderSide.none,
        ),
        focusedBorder: const OutlineInputBorder(
          borderRadius: BorderRadius.all(Radius.circular(12)),
          borderSide: BorderSide(color: Colors.white),
        ),
      ),
      keyboardType: TextInputType.text,
    );
  }
}
