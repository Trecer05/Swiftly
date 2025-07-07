import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/auth/widgets/custom_radio.dart';
import '../../core/themes/theme.dart';

class PrivacyPolicyWidget extends StatelessWidget {
  final bool isChecked;
  final VoidCallback onPressed;
  const PrivacyPolicyWidget({super.key, required this.isChecked, required this.onPressed});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        CustomRadio(isChecked: isChecked, onPressed: onPressed),
        SizedBox(width: 10),
        GestureDetector(
          onTap: () {},
          child: Text(
            'Согласен с политикой конфиденциальности',
            style: AppTextStyles.text_2,
          ),
        ),
      ],
    );
  }
}