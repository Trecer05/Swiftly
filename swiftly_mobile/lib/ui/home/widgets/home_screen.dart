import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/home/widgets/search_field.dart';
import 'package:swiftly_mobile/ui/home/widgets/custom_button.dart';
import 'package:swiftly_mobile/ui/home/widgets/account_widget.dart';

import '../../core/themes/theme.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      backgroundColor: Color.fromARGB(255, 9, 30, 114),
      body: Padding(
        padding: EdgeInsets.symmetric(horizontal: 10, vertical: 20),
        child: Column(
          children: [
            Row(
              children: [
                Text('Главная', style: AppTextStyles.text_6),
                Spacer(),
                Expanded(child: SearchField(hintText: 'Поиск')),
                SizedBox(width: 5),
                CustomButton(),
                SizedBox(width: 5),
                AccountWidget(),
              ],
            )
          ],
        ),
      ));
  }
}
