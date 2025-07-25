import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';

class TextFieldWidget extends StatefulWidget {
  final String hintText;
  final Widget suffixIconWidget;
  final bool isPasswordField;
  const TextFieldWidget({
    super.key,
    required this.hintText,
    required this.suffixIconWidget,
    required this.isPasswordField,
  });

  @override
  State<TextFieldWidget> createState() => _TextFieldWidgetState();
}

class _TextFieldWidgetState extends State<TextFieldWidget> {
  bool obscureText = false;
  void toggleObscureText() {
    setState(() {
      obscureText = !obscureText;
    });
  }

  @override
  Widget build(BuildContext context) {
    return TextField(
      style: AppTextStyles.style4,
      obscureText: obscureText,
      decoration: InputDecoration(
        filled: true,
        fillColor: AppColors.white15,
        hintText: widget.hintText,
        hintStyle: AppTextStyles.style3,
        suffixIcon:
            widget.isPasswordField
                ? VisibilityWidget(
                  isVisible: !obscureText,
                  onPressed: toggleObscureText,
                )
                : widget.suffixIconWidget,
        enabledBorder: const OutlineInputBorder(
          borderRadius: BorderRadius.all(Radius.circular(12)),
          borderSide: BorderSide.none,
        ),
        focusedBorder: const OutlineInputBorder(
          borderRadius: BorderRadius.all(Radius.circular(12)),
          borderSide: BorderSide(color: Colors.white),
        ),
        // disabledBorder: OutlineInputBorder(
        //   borderRadius: BorderRadius.all(Radius.circular(10)),
        //   borderSide: BorderSide(color: AppColors.blue),
        // ),
      ),
      keyboardType: TextInputType.phone,
    );
  }
}

class CheckFillWidget extends StatelessWidget {
  final bool ok;
  const CheckFillWidget({super.key, required this.ok});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(top: 10, right: 10, bottom: 10, left: 10),
      child: Container(
        width: 12,
        height: 12,
        decoration: BoxDecoration(
          color: ok ? AppColors.green : AppColors.red,
          borderRadius: BorderRadius.circular(10),
        ),
        child:
            ok
                ? const Icon(Icons.check, color: AppColors.white, size: 14)
                : const Icon(Icons.clear, color: AppColors.white, size: 14),
      ),
    );
  }
}

class VisibilityWidget extends StatelessWidget {
  final bool isVisible;
  final VoidCallback onPressed;

  const VisibilityWidget({
    super.key,
    required this.isVisible,
    required this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onPressed,
      child: Padding(
        padding: const EdgeInsets.only(
          top: 10,
          right: 10,
          bottom: 10,
          left: 10,
        ),
        child: Icon(
          isVisible ? Icons.visibility : Icons.visibility_off,
          color: AppColors.white,
          size: 20,
        ),
      ),
    );
  }
}
