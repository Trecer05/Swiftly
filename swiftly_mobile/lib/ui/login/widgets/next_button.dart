import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';
import 'privacy_policy_widget.dart';

class NextButton extends StatefulWidget {
  final String buttonText;
  const NextButton({super.key, required this.buttonText});

  @override
  State<NextButton> createState() => _NextButtonState();
}

class _NextButtonState extends State<NextButton> {
  bool policyAccepted = false;

  void toogleexcludeFromSemantics() {
    setState(() {
      policyAccepted = !policyAccepted;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        PrivacyPolicyWidget(
          isChecked: policyAccepted,
          onPressed: toogleexcludeFromSemantics,
        ),
        SizedBox(height: 12),
        GestureDetector(
          onTap: policyAccepted ? () {print('yea!');} : null,
          child: Container(
            width: double.infinity,
            decoration: BoxDecoration(
              color: policyAccepted ? null : AppColors.grey,
              gradient: policyAccepted ? AppColors.gradient_1 : null,
              borderRadius: BorderRadius.circular(12),
            ),
            padding: EdgeInsets.symmetric(vertical: 10),
            child: Center(
              child: Text(widget.buttonText, style: AppTextStyles.text_2),
            ),
          ),
        ),
      ],
    );
  }
}
