import 'package:flutter/material.dart';
import 'package:swiftly_mobile/domain/kanban/models/card_item.dart';

import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';

class DateWidget extends StatelessWidget {
  final DateTime date;
  const DateWidget({super.key, required this.date});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(5),
      decoration: BoxDecoration(
        color: AppColors.white15,
        borderRadius: BorderRadius.circular(15),
      ),
      child: Text(date.toShortDate(), style: AppTextStyles.style9),
    );
  }
}

