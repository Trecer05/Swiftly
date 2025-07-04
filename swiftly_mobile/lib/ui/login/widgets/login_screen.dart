import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/login/widgets/left_panel.dart';
import 'package:swiftly_mobile/ui/login/widgets/right_panel.dart';

import '../../core/themes/colors.dart';

class LoginScreen extends StatelessWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: AppColors.transparent,
        body: Row(
          children: [Expanded(child: LeftPanel()), Expanded(child: RightPanel())],
        ),
      );
  }
}
