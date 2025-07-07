import 'package:flutter/material.dart';

import '../../auth/widgets/next_button.dart';
import '../../auth/widgets/text_field_widget.dart';
import '../../core/themes/theme.dart';

class VerifyCodeScreen extends StatelessWidget {
  const VerifyCodeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
        children: [
          Text('Введите код из письма', style: AppTextStyles.text_1),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                'Код придет на alesha@gmail.com',
                style: AppTextStyles.text_2,
              ),
              Text(' изменить', style: AppTextStyles.text_5),
            ],
          ),
          SizedBox(height: 12),
          TextFieldWidget(
            hintText: 'Код',
            suffixIconWidget: CheckFillWidget(ok: true),
            isPasswordField: false,
          ),
          SizedBox(height: 12),
          Text(
            'Получить новый код можно через 00:30',
            style: AppTextStyles.text_2,
          ),
          SizedBox(height: 12),
          NextButton(buttonText: 'Продолжить', pathScreen: '/home'),
        ],
    );
  }
}
