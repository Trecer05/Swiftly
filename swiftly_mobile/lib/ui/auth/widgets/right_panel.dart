import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';

class RightPanel extends StatelessWidget {
  const RightPanel({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.only(left: 16),
      decoration: BoxDecoration(
        // gradient: AppColors.gradient_2
        color: const Color.fromARGB(255, 64, 64, 64)
      ),
      child: SvgPicture.asset('assets/test.svg'),
    );
  }
}
