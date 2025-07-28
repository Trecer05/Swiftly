import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:swiftly_mobile/ui/auth/widgets/left_panel.dart';
import 'package:swiftly_mobile/ui/auth/widgets/right_panel.dart';

import '../../core/themes/colors.dart';

class AuthScreen extends StatelessWidget {
  const AuthScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.transparent,
      body: Stack(
        children: [
          ImageFiltered(
            imageFilter: ImageFilter.blur(sigmaX: 150, sigmaY: 150),
            child: SvgPicture.asset('assets/test.svg'),
          ),
          const Row(
            children: [
              Expanded(child: LeftPanel()),
              Expanded(child: RightPanel()),
            ],
          ),
        ],
      ),
    );
  }
}
