import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';

class AccountWidget extends StatelessWidget {
  const AccountWidget({super.key});

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
        child: const Icon(Icons.account_circle_outlined, color: AppColors.white128),
      ),
    );
  }
}
