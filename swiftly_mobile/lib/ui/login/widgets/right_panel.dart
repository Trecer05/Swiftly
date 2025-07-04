import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/core/themes/colors.dart';
// import 'package:flutter_svg/svg.dart';

class RightPanel extends StatelessWidget {
  const RightPanel({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(gradient: AppColors.gradient_2),
      child: Padding(
        padding: EdgeInsets.only(left: 16, top: 100, bottom: 100),
        child: SizedBox(
          width: double.infinity,
          height: double.infinity,
          child: Image.asset('assets/board.png'),
        ),
      ),
    );
  }
}
