import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/custom_app_bar_desktop.dart';
import 'package:swiftly_mobile/ui/home/widgets/search_field.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/custom_button.dart';

import '../../core/themes/colors.dart';
import 'content/content.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.transparent168,
      body: Column(
        children: [
          CustomAppBarDesktop(
            title: 'Главная',
            buttons: [
              const Expanded(child: SearchField(hintText: 'Поиск')),
              CustomButton(
                prefixIcon: Icons.add,
                text: 'Пригласить',
                gradient: false,
                onTap: () {},
              ),
              CustomButton(
                prefixIcon: Icons.account_circle_outlined,
                gradient: false,
                onTap: () {},
              ),
            ],
          ),
          const Expanded(
            child: SingleChildScrollView(
              child: Content(),
            ),
          ),
        ],
      ),
    );
  }
}
