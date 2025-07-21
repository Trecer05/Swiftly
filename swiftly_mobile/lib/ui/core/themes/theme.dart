import 'package:flutter/material.dart';
import 'colors.dart';

abstract class AppFontWeights {
  static const FontWeight bolt = FontWeight.bold;
  static const FontWeight bolt500 = FontWeight.w500;
  static const FontWeight bolt400 = FontWeight.w400;
  static const FontWeight bolt600 = FontWeight.w600;
}

abstract class AppFontSizes {
  static const double size40 = 40;
  static const double size32 = 32;
  static const double size26 = 26;
  static const double size22 = 22;
  static const double size20 = 20;
  static const double size16 = 16;
  static const double size14 = 14;
  static const double size12 = 12;
}

abstract class AppFontStyles {}

abstract class AppTextStyles {
  static const TextStyle style1 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.size20,
    fontWeight: AppFontWeights.bolt,
  );
  static const TextStyle style2 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.size12,
    fontWeight: AppFontWeights.bolt500,
  );
  static const TextStyle style3 = TextStyle(
    color: AppColors.white128,
    fontSize: AppFontSizes.size12,
    fontWeight: AppFontWeights.bolt400,
  );
  static const TextStyle style4 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.size12,
    fontWeight: AppFontWeights.bolt400,
  );
  static const TextStyle style5 = TextStyle(
    color: Color.fromARGB(255, 36, 29, 244),
    fontSize: AppFontSizes.size12,
    fontWeight: AppFontWeights.bolt500,
  );
  static const TextStyle style6 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.size32,
    fontWeight: AppFontWeights.bolt,
  );
  static const TextStyle style7 = TextStyle(
    color: AppColors.white128,
    fontSize: AppFontSizes.size16,
    fontWeight: AppFontWeights.bolt500,
  );
  static const TextStyle style8 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.size16,
    fontWeight: AppFontWeights.bolt600,
  );
  static const TextStyle style9 = TextStyle(
    color: AppColors.white128,
    fontSize: AppFontSizes.size12,
    fontWeight: AppFontWeights.bolt500,
  );
  static const TextStyle style10 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.size20,
    fontWeight: AppFontWeights.bolt600,
  );
  static const TextStyle style11 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.size40,
    fontWeight: AppFontWeights.bolt600,
  );
  static const TextStyle style12 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.size20,
    fontWeight: AppFontWeights.bolt500,
  );
  static const TextStyle style13 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.size14,
    fontWeight: AppFontWeights.bolt500,
  );
}
