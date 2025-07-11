import 'package:flutter/material.dart';

import '../../core/themes/colors.dart';
import '../../core/themes/theme.dart';

enum Priority {
  low('Низкий', AppColors.green1, AppColors.green2),
  medium('Средний', AppColors.yellow1, AppColors.yellow2),
  high('Высокий', AppColors.red1, AppColors.red2);

  final String name;
  final Color textColor;
  final Color backgroundColor;

  const Priority(this.name, this.textColor, this.backgroundColor);
}

class CartWidget extends StatelessWidget {
  final CategoryWidget categoryWidget;
  final String name;
  final String description;
  final String data;
  final Priority priority;
  const CartWidget({
    super.key,
    required this.categoryWidget,
    required this.name,
    required this.description,
    required this.data,
    required this.priority,
  });

  @override
  Widget build(BuildContext context) {
    return IntrinsicWidth(
      child: 
      Container(
        padding: const EdgeInsets.all(10),
        decoration: BoxDecoration(
          color: AppColors.white15,
          borderRadius: BorderRadius.circular(15),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                categoryWidget,
                const Spacer(),
                const DataWidget(data: '5 aug'),
              ],
            ),
            const SizedBox(height: 5),
            Text(name, style: AppTextStyles.text_8),
            const SizedBox(height: 5),
            Text(description, style: AppTextStyles.text_9),
            const SizedBox(height: 5),
            PriorityWidget(priority: priority),
          ],
        ),
      ),
    );
  }
}

class DataWidget extends StatelessWidget {
  final String data;
  const DataWidget({super.key, required this.data});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(5),
      decoration: BoxDecoration(
        color: AppColors.white15,
        borderRadius: BorderRadius.circular(15),
      ),
      child: Text(data, style: AppTextStyles.text_9),
    );
  }
}

class CategoryWidget extends StatelessWidget {
  final String name;
  const CategoryWidget({super.key, required this.name});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(5),
      decoration: BoxDecoration(
        color: AppColors.wanderingThrus2,
        borderRadius: BorderRadius.circular(15),
      ),
      child: Text(name, style: AppTextStyles.text_10),
    );
  }
}

class PriorityWidget extends StatelessWidget {
  final Priority priority;
  const PriorityWidget({super.key, required this.priority});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(5),
      decoration: BoxDecoration(
        color: priority.backgroundColor,
        borderRadius: BorderRadius.circular(15),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            priority.name,
            style: TextStyle(
              color: priority.textColor,
              fontSize: AppFontSizes.hm,
              fontWeight: AppFontWeights.bolt_1,
            ),
          ),
          const SizedBox(width: 5),
          Icon(Icons.network_cell_outlined, color: priority.textColor, size: 12,),
        ],
      ),
    );
  }
}
