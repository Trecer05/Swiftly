import 'package:flutter/material.dart';

abstract final class AppColors {
  static const Color white = Color.fromARGB(255, 255, 255, 255);
  static const Color black = Color.fromARGB(255, 0, 0, 0);

  static const Color white15 = Color.fromARGB(15, 255, 255, 255);
  static const Color white128 = Color.fromARGB(128, 255, 255, 255);

  static const Color red = Colors.red;
  static const Color green = Colors.green;
  static const Color blue = Color.fromARGB(255, 19, 92, 236);

  static const Color grey = Colors.grey;

  static const Color transparent = Colors.transparent;

  static const Color background = Color.fromARGB(190, 16, 25, 131);

  static const LinearGradient gradient_1 = LinearGradient(
    colors: [Color.fromRGBO(19, 92, 236, 1), Color.fromRGBO(11, 55, 142, 1)],
    begin: Alignment.centerRight,
    end: Alignment.centerLeft,
  );
  static const LinearGradient gradient_2 = LinearGradient(
    colors: [Color.fromRGBO(19, 92, 236, 1), Color.fromRGBO(11, 55, 142, 1)],
    begin: Alignment.bottomCenter,
    end: Alignment.topCenter,
  );
  static const LinearGradient gradient_3 = LinearGradient(
    colors: [Color.fromRGBO(19, 92, 236, 1), Color.fromRGBO(11, 55, 142, 1)],
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
  );
}
