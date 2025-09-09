import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../providers/card_notifier_provider.dart';
import '../../../../providers/label_notifier_provider.dart';
import '../../themes/colors.dart';
import '../../../../domain/kanban/models/card_item.dart';
import '../../../kanban/widgets/card_details_widgets/cart_details_widget.dart';
import 'card_item_content.dart';

class CardItemMobile extends ConsumerWidget {
  final CardItem card;
  const CardItemMobile({super.key, required this.card});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final labels =
        ref
            .watch(labelNotifierProvider)
            .labels
            .where((label) => label.cardId == card.id)
            .toList();

    return GestureDetector(
      onTap: () => _openDialog(context),
      onLongPressStart: (details) => _openMenu(context, details, ref),
      child: Container(
        width: double.maxFinite,
        padding: const EdgeInsets.all(10),
        decoration: BoxDecoration(
          color: AppColors.white15,
          borderRadius: BorderRadius.circular(15),
        ),
        child: CardItemContent(card: card, labels: labels),
      ),
    );
  }

  void _openDialog(BuildContext context) {
    showDialog(
      context: context,
      builder:
          (context) => Dialog(
            insetPadding: const EdgeInsets.all(16),
            backgroundColor: Colors.transparent,
            child: CartDetailsWidget(card: card),
          ),
    );
  }

  void _openMenu(
    BuildContext context,
    LongPressStartDetails details,
    WidgetRef ref,
  ) {
    showMenu(
      context: context,
      position: RelativeRect.fromLTRB(
        details.globalPosition.dx,
        details.globalPosition.dy,
        details.globalPosition.dx,
        details.globalPosition.dy,
      ),
      color: Colors.grey,
      items: [
        PopupMenuItem(
          child: const Text('Удалить'),
          onTap: () {
            ref.read(cardNotifierProvider.notifier).removeCard(card.id);
          },
        ),
        PopupMenuItem(
          child: const Text('Редактировать'),
          onTap: () => _openDialog(context),
        ),
      ],
    );
  }
}
