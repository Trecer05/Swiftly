import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';

class CallMadeWidget extends StatelessWidget {
  const CallMadeWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 32,
      width: 32,
      decoration: BoxDecoration(
        color: AppColors.white128,
        borderRadius: BorderRadius.circular(18.5)
      ),
      child: const Icon(Icons.call_made, color: AppColors.white),
    );
  }
}