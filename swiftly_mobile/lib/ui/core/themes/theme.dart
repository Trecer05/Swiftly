import 'package:flutter/material.dart';
import 'colors.dart';

abstract class AppFontWeights {
  static const FontWeight bolt = FontWeight.bold;
  static const FontWeight bolt_1 = FontWeight.w500;
  static const FontWeight bolt_2 = FontWeight.w400;
}

abstract class AppFontSizes {
  static const double impressive = 36;
  static const double big = 26;
  static const double giant = 22;
  static const double large = 20;
  static const double medium = 16;
  static const double small = 14;
  static const double hm = 12;
}

abstract class AppFontStyles {}

abstract class AppTextStyles {
  static const TextStyle text_1 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.large,
    fontWeight: AppFontWeights.bolt,
  );
  static const TextStyle text_2 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.hm,
    fontWeight: AppFontWeights.bolt_1,
  );
  static const TextStyle text_3 = TextStyle(
    color: AppColors.white128,
    fontSize: AppFontSizes.hm,
    fontWeight: AppFontWeights.bolt_2,
  );
  static const TextStyle text_4 = TextStyle(
    color: AppColors.white,
    fontSize: AppFontSizes.hm,
    fontWeight: AppFontWeights.bolt_2,
  );
  static const TextStyle text_5 = TextStyle(
    color: Color.fromARGB(255, 36, 29, 244),
    fontSize: AppFontSizes.hm,
    fontWeight: AppFontWeights.bolt_1,
  );
}
