import 'dart:ui';

import '../../../ui/core/themes/colors.dart';

enum Priority { 
  low('Low', 'Низкий', AppColors.green), 
  medium('Medium', 'Средний', AppColors.yellow),
  high('High', 'Высокий', AppColors.red); 

  final String id;
  final String title;
  final Color color;

  const Priority(this.id, this.title, this.color);
}