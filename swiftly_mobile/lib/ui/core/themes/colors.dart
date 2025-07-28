import 'package:flutter/material.dart';

abstract final class AppColors {
  static const Color white = Color.fromARGB(255, 255, 255, 255);
  static const Color black = Color.fromARGB(255, 0, 0, 0);

  static const Color black168 = Color.fromARGB(168, 8, 8, 8);
  static const Color white15 = Color.fromARGB(15, 255, 255, 255);
  static const Color white26 = Color.fromARGB(26, 255, 255, 255);
  static const Color white38 = Color.fromARGB(38, 255, 255, 255);
  static const Color white64 = Color.fromARGB(64, 255, 255, 255);
  static const Color white128 = Color.fromARGB(128, 255, 255, 255);

  static const Color red = Colors.red;
  static const Color green = Colors.green;
  static const Color blue = Color.fromARGB(255, 19, 92, 236);
  static const Color yellow = Colors.yellow;
  static const Color grey = Colors.grey;

  static const Color wanderingThrus = Color.fromARGB(255, 31, 217, 186);
  static const Color amaranthMagenta = Color.fromARGB(255, 217, 31, 192);
  static const Color twitter = Color.fromARGB(74, 52, 148, 231);
  static const Color green1 = Color.fromARGB(255, 31, 217, 72);
  static const Color green2 = Color.fromARGB(51, 31, 217, 72);
  static const Color yellow1 = Color.fromARGB(255, 217, 167, 31);
  static const Color yellow2 = Color.fromARGB(51, 217, 167, 31);
  static const Color red1 = Color.fromARGB(255, 217, 31, 31);
  static const Color red2 = Color.fromARGB(51, 217, 31, 31);

  static const Color transparent = Color.fromARGB(0, 0, 0, 0);
  static const Color transparent168 = Color.fromARGB(168, 8, 8, 8);
  static const Color transparent191 = Color.fromARGB(191, 8, 8, 8);

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
  static const LinearGradient gradient_4 = LinearGradient(
    colors: [Color.fromRGBO(38, 98, 217, 1), Color.fromRGBO(23, 59, 130, 1)],
    begin: Alignment.centerRight,
    end: Alignment.centerLeft,
  );
}
