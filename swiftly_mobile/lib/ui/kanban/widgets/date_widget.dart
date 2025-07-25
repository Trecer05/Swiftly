import 'package:flutter/material.dart';
import 'package:swiftly_mobile/domain/kanban/models/card_item.dart';

import '../../core/themes/theme.dart';

class DateWidget extends StatelessWidget {
  final DateTime date;
  final Color color;
  const DateWidget({super.key, required this.date, required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(5),
      decoration: BoxDecoration(
        color: color,
        borderRadius: BorderRadius.circular(15),
      ),
      child: Text(date.toShortDate(), style: AppTextStyles.style9),
    );
  }
}

