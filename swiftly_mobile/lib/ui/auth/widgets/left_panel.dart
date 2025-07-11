import 'package:flutter/material.dart';

import '../../register/widgets/register_screen.dart';

class LeftPanel extends StatefulWidget {
  const LeftPanel({super.key});

  @override
  State<LeftPanel> createState() => _LeftPanelState();
}

class _LeftPanelState extends State<LeftPanel> {
  late Widget currentScreen;

  @override
  void initState() {
    super.initState();
    currentScreen = RegisterScreen(
      onTap: (newScreen) {
        setState(() {
          currentScreen = newScreen;
        });
      },
    );
  }
  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: const BoxDecoration(
        // gradient: AppColors.gradient_3
        // color: Color.fromARGB(255, 8, 19, 55),
        color: Colors.transparent,
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            currentScreen,
          ],
        ),
      ),
    );
  }
}
