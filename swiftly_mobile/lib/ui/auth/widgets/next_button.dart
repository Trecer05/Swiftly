import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:swiftly_mobile/ui/verify_code/widgets/verify_code_screen.dart';

import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';
import 'privacy_policy_widget.dart';

class NextButton extends StatefulWidget {
  final String buttonText;
  final String? pathScreen;
  final ValueChanged<Widget>? onTap;
  const NextButton({super.key, required this.buttonText, this.pathScreen, this.onTap});

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
        const SizedBox(height: 12),
        GestureDetector(
          onTap: policyAccepted ? () {
          if (widget.onTap != null) {
            widget.onTap!(const VerifyCodeScreen());
          } else if (widget.pathScreen != null) {
            context.go(widget.pathScreen!);
          }
        }
      : null,
          child: Container(
            width: double.infinity,
            decoration: BoxDecoration(
              color: policyAccepted ? null : AppColors.grey,
              gradient: policyAccepted ? AppColors.gradient_1 : null,
              borderRadius: BorderRadius.circular(12),
            ),
            padding: const EdgeInsets.symmetric(vertical: 10),
            child: Center(
              child: Text(widget.buttonText, style: AppTextStyles.text_2),
            ),
          ),
        ),
      ],
    );
  }
}
