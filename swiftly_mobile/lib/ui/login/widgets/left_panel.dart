import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/core/themes/colors.dart';
import 'package:swiftly_mobile/ui/login/widgets/container_widget.dart';
import 'package:swiftly_mobile/ui/login/widgets/next_button.dart';
import 'package:swiftly_mobile/ui/login/widgets/text_field_widget.dart';

import '../../core/themes/theme.dart';

class LeftPanel extends StatelessWidget {
  const LeftPanel({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        gradient: AppColors.gradient_3
      ),
      child: Padding(
        padding: EdgeInsets.symmetric(horizontal: 16),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Text('Создайте аккаунт', style: AppTextStyles.text_1),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text('или', style: AppTextStyles.text_2),
                Text(' войдите', style: AppTextStyles.text_5),
              ],
            ),
            SizedBox(height: 12),
            TextFieldWidget(
              hintText: 'Имя',
              suffixIconWidget: CheckFillWidget(ok: true),
              isPasswordField: false,
            ),
            SizedBox(height: 12),
            TextFieldWidget(
              hintText: 'Почта',
              suffixIconWidget: CheckFillWidget(ok: false),
              isPasswordField: false,
            ),
            SizedBox(height: 12),
            TextFieldWidget(
              hintText: 'Пароль',
              suffixIconWidget: const SizedBox(),
              isPasswordField: true,
            ),
            SizedBox(height: 12),
            NextButton(buttonText: 'Продолжить'),
            SizedBox(height: 12),
            Divider(thickness: 2, color: AppColors.white15),
            SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: ContainerWidget(imagePath: 'assets/google_logo.png'),
                ),
                SizedBox(width: 10),
                Expanded(child: ContainerWidget(imagePath: 'assets/vk_logo.png')),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
