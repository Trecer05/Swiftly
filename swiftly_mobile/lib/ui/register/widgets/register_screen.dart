import 'package:flutter/material.dart';

import '../../auth/widgets/container_widget.dart';
import '../../auth/widgets/next_button.dart';
import '../../auth/widgets/text_field_widget.dart';
import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';

class RegisterScreen extends StatelessWidget {
  final ValueChanged<Widget>? onTap;
  const RegisterScreen({super.key, this.onTap});

  @override
  Widget build(BuildContext context) {
    return Column(
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
        NextButton(buttonText: 'Продолжить', onTap: onTap),
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
    );
  }
}
