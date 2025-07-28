import 'package:flutter/material.dart';
import 'package:swiftly_mobile/ui/home/widgets/content/active_tasks.dart';

import '../../../core/themes/colors.dart';
import '../../../core/themes/theme.dart';
import 'participants.dart';

class Content extends StatelessWidget {
  const Content({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(left: 10, right: 10, top: 20, bottom: 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: 50),
          const Padding(
            padding: EdgeInsets.symmetric(horizontal: 20),
            child: Text('Swiftly project', style: AppTextStyles.style11),
          ),
          const SizedBox(height: 24),
          Container(
            padding: const EdgeInsets.only(
              left: 24,
              right: 24,
              top: 24,
              bottom: 16,
            ),
            decoration: BoxDecoration(
              color: AppColors.white15,
              borderRadius: BorderRadius.circular(12),
            ),
            child: const Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Participants(),
                SizedBox(height: 20),
                Text('Мои активные задачи', style: AppTextStyles.style12),
                SizedBox(height: 10),
                ActiveTasks(),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
