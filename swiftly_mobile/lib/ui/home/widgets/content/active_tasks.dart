import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../providers/card_notifier_provider.dart';
import '../../../../providers/current_user_provider.dart';
import '../../../core/ui/card_item/card_item_desktop.dart';

class ActiveTasks extends ConsumerWidget {
  const ActiveTasks({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final me = ref.watch(currentUserProvider);
    final cards =
        ref
            .watch(cardNotifierProvider)
            .cards
            .where((card) => card.userId == me!.id)
            .toList();
    return Wrap(
      spacing: 10,
      runSpacing: 10,
      children: [
        ...cards.map(
          (card) => CardItemDesktop(
            card: card,
          ),
        ),
      ],
    );
  }
}
