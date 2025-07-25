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
        const Text('Создайте аккаунт', style: AppTextStyles.style1),
        const Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('или', style: AppTextStyles.style2),
            Text(' войдите', style: AppTextStyles.style5),
          ],
        ),
        const SizedBox(height: 12),
        const TextFieldWidget(
          hintText: 'Имя',
          suffixIconWidget: CheckFillWidget(ok: true),
          isPasswordField: false,
        ),
        const SizedBox(height: 12),
        const TextFieldWidget(
          hintText: 'Почта',
          suffixIconWidget: CheckFillWidget(ok: false),
          isPasswordField: false,
        ),
        const SizedBox(height: 12),
        const TextFieldWidget(
          hintText: 'Пароль',
          suffixIconWidget: SizedBox(),
          isPasswordField: true,
        ),
        const SizedBox(height: 12),
        NextButton(buttonText: 'Продолжить', onTap: onTap),
        const SizedBox(height: 12),
        const Divider(thickness: 2, color: AppColors.white15),
        const SizedBox(height: 12),
        const Row(
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
