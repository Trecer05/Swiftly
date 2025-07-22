import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/custom_app_bar.dart';
import 'package:swiftly_mobile/ui/home/widgets/search_field.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/custom_button.dart';

import '../../core/themes/theme.dart';
import 'content/content.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color.fromARGB(255, 9, 30, 114),
      body: Column(
        children: [
          CustomAppBar(
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
              padding: EdgeInsets.only(
                left: 10,
                right: 10,
                top: 20,
                bottom: 20,
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  SizedBox(height: 50),
                  Padding(
                    padding: EdgeInsets.symmetric(horizontal: 20),
                    child: Text(
                      'Swiftly project',
                      style: AppTextStyles.style11,
                    ),
                  ),
                  SizedBox(height: 24),
                  Content(),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
