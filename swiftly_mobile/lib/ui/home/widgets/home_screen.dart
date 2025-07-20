import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/home/widgets/search_field.dart';
import 'package:swiftly_mobile/ui/home/widgets/custom_button.dart';

import '../../core/themes/theme.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color.fromARGB(255, 9, 30, 114),
      body: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Text('Главная', style: AppTextStyles.style6),
                const Spacer(),
                const Expanded(child: SearchField(hintText: 'Поиск')),
                const SizedBox(width: 5),
                CustomButton(prefixIcon: Icons.add, text: 'Пригласить', onTap: (){},),
                const SizedBox(width: 5),
                CustomButton(prefixIcon: Icons.account_circle_outlined, onTap: (){},),
              ],
            ),
            const SizedBox(height: 50),
            const Padding(padding: EdgeInsets.symmetric(horizontal: 20), child: Text('Swiftly project', style: AppTextStyles.style11,))
          ],
        ),
      ));
  }
}
