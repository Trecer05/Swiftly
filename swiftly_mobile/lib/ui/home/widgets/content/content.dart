import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/home/widgets/content/active_tasks.dart';

import '../../../core/themes/colors.dart';
import '../../../core/themes/theme.dart';
import 'participants.dart';

class Content extends StatelessWidget {
  const Content({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
        // width: double.infinity,
        padding: const EdgeInsets.only(left: 24, right: 24, top: 24, bottom: 16),
        decoration: BoxDecoration(
          color: AppColors.white15,
          borderRadius: BorderRadius.circular(12),
        ),
        child: const Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
          Participants(),
          SizedBox(height: 20),
          Text('Мои активные задачи', style: AppTextStyles.style12,),
          SizedBox(height: 10),
          ActiveTasks(),
        ],),
    );
  }
}
