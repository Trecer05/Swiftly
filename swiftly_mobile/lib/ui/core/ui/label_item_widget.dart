import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/providers/label_notifier_provider.dart';
import 'package:swiftly_mobile/ui/core/state/label_item_settings.dart';

import '../../../domain/models/label_item.dart';
import '../themes/theme.dart';

class LabelItemWidget extends StatelessWidget {
  final LabelItem labelItem;

  const LabelItemWidget({super.key, required this.labelItem});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {
        showDialog(
          context: context,
          builder:
              (context) => AlertDialog(
                backgroundColor: Colors.transparent,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                contentPadding: EdgeInsets.zero,
                content: LabelItemSettings(labelItem: labelItem),
              ),
        );
      },
      child: Container(
        padding: const EdgeInsets.all(5),
        decoration: BoxDecoration(
          color: labelItem.color.withValues(alpha: 0.2),
          borderRadius: BorderRadius.circular(15),
        ),
        child: Text(
          labelItem.title,
          style: TextStyle(
            color: labelItem.color,
            fontSize: AppFontSizes.size12,
            fontWeight: AppFontWeights.bolt500,
          ),
        ),
      ),
    );
  }
}

class TestLabel extends ConsumerWidget {
  const TestLabel({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final labels =
        ref
            .watch(labelNotifierProvider)
            .labels
            .toList();

    return ListView.builder(
      itemCount: labels.length,
      itemBuilder: (context, index) {
        final label = labels[index];
        return Padding(
          padding: const EdgeInsets.symmetric(vertical: 4, horizontal: 8),
          child: LabelItemWidget(labelItem: label),
        );
      },
    );
  }
}
